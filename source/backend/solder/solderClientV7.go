package solder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/mitchellh/mapstructure"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/crawlers"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SolderClientV7 struct {
	Client            *http.Client
	Url               *url.URL
	modVersionIdCache map[string]map[string]string
	modIdCache        map[string]string
	buildCache        map[string]crawlers.Build
	lock              sync.RWMutex
	modVersionIdLock  sync.RWMutex
	modListSemaphor   chan struct{}
	addMutex          sync.Mutex
	buildMutex        sync.Mutex
}

func NewV7SolderClient(Url string) *SolderClientV7 {
	var cookieJar, _ = cookiejar.New(nil)

	var client = &http.Client{
		Jar: cookieJar,
	}

	u, err := url.Parse(Url)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Panic(err)
	}

	return &SolderClientV7{
		Client:            client,
		Url:               u,
		modVersionIdCache: make(map[string]map[string]string),
		modIdCache:        make(map[string]string),
		buildCache:        make(map[string]crawlers.Build),
		modListSemaphor:   make(chan struct{}, 20),
	}
}

func TestSolderConnection(conn types.WebsocketConnection, data interface{}) {
	conn.Write("solder-test", "Starting solder test")
	dict := data.(map[string]interface{})

	var solderInfo types.SolderInfo
	err := mapstructure.Decode(dict, &solderInfo)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Panic(err)
	}

	client := NewV7SolderClient(solderInfo.Url)
	err = client.Login(solderInfo.Username, solderInfo.Password)
	if err == nil {
		conn.Write("solder-test", "Solder connection seems alright")
	} else {
		conn.Write("solder-test", "Could not connect to solder and login\n"+err.Error())
	}
}

func (s *SolderClientV7) createUrl(after string) url.URL {
	url := *s.Url
	url.Path = path.Join(url.Path, after)
	log.Println(url.Path)
	return url
}

func (s *SolderClientV7) postForm(url string, data url.Values) *http.Response {
	s.modListSemaphor <- struct{}{}
	defer func() { <-s.modListSemaphor }()

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	if err != nil {
		raven.CaptureError(err, nil)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := s.Client.Do(req)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Panic(err)
	}
	return resp
}

func (s *SolderClientV7) doRequest(method, url, data string) *http.Response {
	s.modListSemaphor <- struct{}{}
	defer func() { <-s.modListSemaphor }()

	var body io.Reader
	if data != "" {
		body = strings.NewReader(data)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		raven.CaptureError(err, nil)
	}
	// Yes, this request is totally an ajax request...
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	response, err := s.Client.Do(req)
	if err != nil {
		raven.CaptureError(err, nil)
		log.Panic(err)
	}
	return response
}

func (s *SolderClientV7) Login(email string, password string) error {
	Url := s.createUrl("login")

	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)

	response := s.postForm(Url.String(), form)
	defer response.Body.Close()

	return crawlers.CrawlLogin(response)
}

func IsOnSolder(client *SolderClientV7, m *types.Mod) bool {
	if client == nil {
		return m.IsOnSolder
	}
	id := client.GetModId(m.ModId)
	if id == "" {
		return false
	}
	id = client.GetModVersionId(m.GenerateSimpleOutputInfo())
	if id == "" {
		return false
	}
	return true
}

func (s *SolderClientV7) CreatePack(name, slug string) string {
	id := s.GetModpackId(slug)
	if id != "" {
		return id
	}
	Url := s.createUrl("modpack/create")

	form := url.Values{}
	form.Add("name", name)
	form.Add("slug", slug)
	form.Add("hidden", "false")
	res := s.postForm(Url.String(), form)
	defer res.Body.Close()

	u := res.Request.URL
	segments := strings.Split(u.Path, "/")
	segment := segments[len(segments)-1]
	if _, err := strconv.Atoi(segment); err == nil {
		return segment
	} else {
		return s.GetModpackId(slug)
	}

}

func (s *SolderClientV7) AddMod(mod *types.OutputInfo) string {
	// Ensure we don't attempt to create the mod multiple times, as this won't work
	s.addMutex.Lock()
	defer s.addMutex.Unlock()

	modId := s.GetModId(mod.Id)
	if modId != "" {
		return modId
	}

	Url := s.createUrl("mod/create")

	form := url.Values{}
	form.Add("pretty_name", mod.Name)
	form.Add("name", mod.Id)
	form.Add("author", mod.Author)
	form.Add("description", mod.Description)
	form.Add("link", mod.Url)
	log.Printf("Adding mod:\n%v\n", form)

	response := s.postForm(Url.String(), form)
	defer response.Body.Close()

	modResponse := crawlers.CrawlMod(response)
	log.Println(modResponse)
	return modResponse.Id
}

func (s *SolderClientV7) AddModVersion(modId, md5, version string) {
	form := url.Values{}
	form.Add("mod-id", modId)
	form.Add("add-version", version)
	form.Add("add-md5", md5)
	Url := s.createUrl("mod/add-version")
	log.Printf("Adding mod version:\n%v\n", form)

	res := s.postForm(Url.String(), form)
	defer res.Body.Close()
}

func (s *SolderClientV7) RehashModVersion(modversionId string, md5 string) {
	form := url.Values{}
	form.Add("version-id", modversionId)
	form.Add("md5", md5)
	Url := s.createUrl("mod/rehash")
	res := s.postForm(Url.String(), form)
	res.Body.Close()
}

func (s *SolderClientV7) IsModversionOnline(mod *types.OutputInfo) bool {
	return s.GetModVersionId(mod) != ""
}

func (s *SolderClientV7) GetBuild(buildId string) crawlers.Build {
	s.buildMutex.Lock()
	defer s.buildMutex.Unlock()

	s.lock.RLock()
	build, ok := s.buildCache[buildId]
	s.lock.RUnlock()
	if ok {
		return build
	}

	Url := s.createUrl("modpack/build/" + buildId)
	res := s.doRequest(http.MethodGet, Url.String(), "")
	defer res.Body.Close()
	build = crawlers.CrawlBuild(res)

	s.lock.Lock()
	s.buildCache[buildId] = build
	s.lock.Unlock()
	return build
}

func (s *SolderClientV7) IsModInBuild(mod *types.OutputInfo, buildId string) bool {
	build := s.GetBuild(buildId)

	for _, m := range build.Mods {
		//log.Printf("%s === %s\n", m.Name, mod.Id)
		if m.Name == mod.Id {
			return true
		}
	}
	return false
}

func (s *SolderClientV7) IsModversionActiveInBuild(mod *types.OutputInfo, buildId string) bool {
	build := s.GetBuild(buildId)

	for _, m := range build.Mods {
		if m.Name == mod.Id && m.Active == mod.GenerateOnlineVersion() {
			return true
		}
	}
	return false
}

func (s *SolderClientV7) GetActiveModversionInBuildId(mod *types.OutputInfo, buildId string) string {
	Url := s.createUrl("modpack/build/" + buildId)
	res := s.doRequest(http.MethodGet, Url.String(), "")
	defer res.Body.Close()
	build := crawlers.CrawlBuild(res)

	var version string
	for _, m := range build.Mods {
		if m.Name == mod.Id {
			version = m.Active
			break
		}
	}
	if version == "" {
		return ""
	}
	return s.GetModVersionId(mod)
}

func (s *SolderClientV7) SetModVersionInBuild(mod *types.OutputInfo, buildId string) {
	modVersionId := s.GetModVersionId(mod)
	Url := s.createUrl("modpack/build/modify")

	form := make(map[string]interface{})
	form["action"] = "version"
	form["build-id"] = buildId
	form["version"] = modVersionId
	form["modversion-id"] = s.GetActiveModversionInBuildId(mod, buildId)
	data, err := json.Marshal(form)
	if err != nil {
		raven.CaptureError(err, nil)
	}

	res := s.doRequest(http.MethodPost, Url.String(), string(data))
	res.Body.Close()
}

func (s *SolderClientV7) IsPackOnline(modpack *types.Modpack) bool {
	return s.GetModpackId(modpack.GetSlug()) != ""
}

func (s *SolderClientV7) IsBuildOnline(modpack *types.Modpack) bool {
	id := s.GetBuildId(modpack)
	return id != ""
}

func (s *SolderClientV7) GetModVersionId(mod *types.OutputInfo) string {
	modId, modVersion := mod.Id, mod.GenerateOnlineVersion()
	s.modVersionIdLock.RLock()
	l, ok := s.modVersionIdCache[modId]
	s.modVersionIdLock.RUnlock()
	if ok {
		id, ok := l[modVersion]
		if ok {
			return id
		}
	}
	solderModId := s.GetModId(modId)
	Url := s.createUrl("mod/view/" + solderModId)
	res := s.doRequest(http.MethodGet, Url.String(), "")
	defer res.Body.Close()
	modVersions := crawlers.CrawlModVersion(res)
	var id string

	for _, mv := range modVersions {
		if mv.Version == modVersion {
			id = mv.Id
			break
		}
	}

	if id != "" {
		s.modVersionIdLock.Lock()
		l, contains := s.modVersionIdCache[modId]
		if !contains {
			l = make(map[string]string, 0)
			s.modVersionIdCache[modId] = l
		}
		l[modVersion] = id
		s.modVersionIdLock.Unlock()
	}

	return id
}

func (s *SolderClientV7) CreateBuild(modpack *types.Modpack, modpackId string) string {
	Url := s.createUrl("modpack/add-build/" + modpackId)

	form := url.Values{}
	form.Add("version", modpack.Version)
	form.Add("minecraft", modpack.MinecraftVersion)
	form.Add("java-version", modpack.Technic.Java)
	form.Add("memory-enabled", strconv.FormatBool(modpack.Technic.Memory != 0))
	form.Add("memory", strconv.FormatInt(modpack.Technic.Memory, 10))
	res := s.postForm(Url.String(), form)
	defer res.Body.Close()

	// Return the build id because the response redirects
	u := res.Request.URL
	segments := strings.Split(u.Path, "/")
	return segments[len(segments)-1]
}

func (s *SolderClientV7) GetBuildId(modpack *types.Modpack) string {
	modpackId := s.GetModpackId(modpack.GetSlug())
	Url := s.createUrl("modpack/view/" + modpackId)

	res := s.doRequest(http.MethodGet, Url.String(), "")
	builds := crawlers.CrawlBuildList(res)

	for _, build := range builds {
		log.Println(build.Version)
		if build.Version == modpack.Version {
			return build.Id
		}
	}
	return ""
}

func (s *SolderClientV7) GetModpackId(slug string) string {
	Url := s.createUrl("modpack/list")

	res := s.doRequest(http.MethodGet, Url.String(), "")
	defer res.Body.Close()
	modpacks := crawlers.CrawlModpackList(res)

	for _, modpack := range modpacks {
		if modpack.Name == slug {
			return modpack.Id
		}
	}
	return ""
}

func (s *SolderClientV7) GetModId(modid string) string {
	s.lock.RLock()
	if id, ok := s.modIdCache[modid]; ok {
		s.lock.RUnlock()
		if id != "" {

			return id
		}
	}
	s.lock.RUnlock()

	Url := s.createUrl("mod/list")
	Url.Query().Add("bust", strconv.FormatInt(time.Now().UnixNano(), 10))

	response := s.doRequest(http.MethodGet, Url.String(), "")
	defer response.Body.Close()
	mods := crawlers.CrawlModList(response)
	for _, mod := range mods {
		id := mod.Id
		if strings.Trim(id, " ") != "" {
			s.lock.Lock()
			s.modIdCache[mod.Name] = mod.Id
			s.lock.Unlock()
		}
		if strings.ToLower(mod.Name) == strings.ToLower(modid) {
			return id
		}
	}
	fmt.Println("Could not find match for " + modid)
	return ""
}

type addBuildToModpackResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func (s *SolderClientV7) AddModversionToBuild(mod *types.OutputInfo, modpackBuildId string) {
	Url := s.createUrl("modpack/modify/add")

	form := url.Values{}
	form.Add("build", modpackBuildId)
	form.Add("mod-name", mod.Id)
	form.Add("mod-version", mod.GenerateOnlineVersion())
	form.Add("action", "add")

	response := s.postForm(Url.String(), form)
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	var res addBuildToModpackResponse
	err := json.NewDecoder(buf).Decode(&res)

	if err != nil {
		log.Println(Url.String())
		log.Printf("%v", form)
		log.Println(buf.String())
		raven.CaptureErrorAndWait(err, nil)
		log.Panic(err)
	}

	if res.Status != "success" {
		if res.Reason == "Duplicate Modversion found" {
			log.Println(*mod, "was already added to the build. For some reason")
		} else {
			log.Println(res.Reason)
			log.Println(*mod)
			log.Panic("Something went wrong when adding a mod to a build, see above mod details")
		}
	}

}
