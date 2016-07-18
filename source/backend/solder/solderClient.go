package solder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/crawlers"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"github.com/zlepper/go-websocket-connection"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type SolderClient struct {
	Client            *http.Client
	Url               *url.URL
	modVersionIdCache map[string]map[string]string
	modIdCache        map[string]string
	buildCache        map[string]crawlers.Build
}

func NewSolderClient(Url string) *SolderClient {
	var cookieJar, _ = cookiejar.New(nil)

	var client = &http.Client{
		Jar: cookieJar,
	}

	u, err := url.Parse(Url)
	if err != nil {
		log.Panic(err)
	}

	return &SolderClient{
		Client:            client,
		Url:               u,
		modVersionIdCache: make(map[string]map[string]string),
		modIdCache:        make(map[string]string),
		buildCache:        make(map[string]crawlers.Build),
	}
}

func TestSolderConnection(conn websocket.WebsocketConnection, data interface{}) {
	dict := data.(map[string]interface{})

	var solderInfo types.SolderInfo
	err := mapstructure.Decode(dict, &solderInfo)
	if err != nil {
		log.Panic(err)
	}

	client := NewSolderClient(solderInfo.Url)
	loginSuccess := client.Login(solderInfo.Username, solderInfo.Password)
	if loginSuccess {
		conn.Write("solder-test", "TECHNIC.SOLDER.SUCCESS")
	} else {
		conn.Write("solder-test", "TECHNIC.SOLDER.ERROR")
	}
}

func (s *SolderClient) createUrl(after string) url.URL {
	url := *s.Url
	url.Path = path.Join(url.Path, after)
	return url
}

func (s *SolderClient) postForm(url string, data url.Values) *http.Response {
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := s.Client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	return resp
}

func (s *SolderClient) doRequest(method, url, data string) *http.Response {
	var body io.Reader
	if data != "" {
		body = strings.NewReader(data)
	}
	req, _ := http.NewRequest(method, url, body)
	// Yes, this request is totally an ajax request...
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	response, err := s.Client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	return response
}

func (s *SolderClient) Login(email string, password string) bool {
	Url := s.createUrl("login")

	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)

	response := s.postForm(Url.String(), form)
	defer response.Body.Close()

	return crawlers.CrawlLogin(response)
}

func (s *SolderClient) CreatePack(name, slug string) string {
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

func (s *SolderClient) AddMod(mod *types.OutputInfo) string {
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
	return s.GetModId(mod.Id)
}

func (s *SolderClient) AddModVersion(modId, md5, version string) {
	form := url.Values{}
	form.Add("mod-id", modId)
	form.Add("add-version", version)
	form.Add("add-md5", md5)
	Url := s.createUrl("mod/add-version")
	log.Printf("Adding mod version:\n%v\n", form)

	res := s.postForm(Url.String(), form)
	defer res.Body.Close()
}

func (s *SolderClient) RehashModVersion(modversionId string, md5 string) {
	form := url.Values{}
	form.Add("version-id", modversionId)
	form.Add("md5", md5)
	Url := s.createUrl("mod/rehash")
	res := s.postForm(Url.String(), form)
	res.Body.Close()
}

func (s *SolderClient) IsModversionOnline(mod *types.OutputInfo) bool {
	return s.GetModVersionId(mod) != ""
}

func (s *SolderClient) GetBuild(buildId string) crawlers.Build {
	if build, ok := s.buildCache[buildId]; ok {
		return build
	}

	Url := s.createUrl("modpack/build/" + buildId)
	res := s.doRequest(http.MethodGet, Url.String(), "")
	defer res.Body.Close()
	build := crawlers.CrawlBuild(res)

	s.buildCache[buildId] = build
	return build
}

func (s *SolderClient) IsModInBuild(mod *types.OutputInfo, buildId string) bool {
	build := s.GetBuild(buildId)

	for _, m := range build.Mods {
		//log.Printf("%s === %s\n", m.Name, mod.Id)
		if m.Name == mod.Id {
			return true
		}
	}
	return false
}

func (s *SolderClient) IsModversionActiveInBuild(mod *types.OutputInfo, buildId string) bool {
	build := s.GetBuild(buildId)

	for _, m := range build.Mods {
		if m.Name == mod.Id && m.Active == mod.GenerateOnlineVersion() {
			return true
		}
	}
	return false
}

func (s *SolderClient) GetActiveModversionInBuildId(mod *types.OutputInfo, buildId string) string {
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

func (s *SolderClient) SetModVersionInBuild(mod *types.OutputInfo, buildId string) {
	modVersionId := s.GetModVersionId(mod)
	Url := s.createUrl("modpack/build/modify")

	form := make(map[string]interface{})
	form["action"] = "version"
	form["build-id"] = buildId
	form["version"] = modVersionId
	form["modversion-id"] = s.GetActiveModversionInBuildId(mod, buildId)
	data, _ := json.Marshal(form)

	res := s.doRequest(http.MethodPost, Url.String(), string(data))
	res.Body.Close()
}

func (s *SolderClient) IsPackOnline(modpack *types.Modpack) bool {
	return s.GetModpackId(modpack.GetSlug()) != ""
}

func (s *SolderClient) IsBuildOnline(modpack *types.Modpack) bool {
	id := s.GetBuildId(modpack)
	return id != ""
}

func (s *SolderClient) GetModVersionId(mod *types.OutputInfo) string {
	modId, modVersion := mod.Id, mod.GenerateOnlineVersion()
	l, ok := s.modVersionIdCache[modId]
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
		l, contains := s.modVersionIdCache[modId]
		if !contains {
			l = make(map[string]string, 0)
			s.modVersionIdCache[modId] = l
		}
		l[modVersion] = id
	}

	return id
}

func (s *SolderClient) CreateBuild(modpack *types.Modpack, modpackId string) string {
	Url := s.createUrl("modpack/add-build/" + modpackId)

	form := url.Values{}
	form.Add("version", modpack.Version)
	form.Add("minecraft", modpack.MinecraftVersion)
	form.Add("java-version", modpack.Technic.Java)
	form.Add("memory-enabled", strconv.FormatBool(modpack.Technic.Memory != 0))
	form.Add("memory", strconv.FormatFloat(modpack.Technic.Memory, 'E', -1, 64))
	res := s.postForm(Url.String(), form)
	defer res.Body.Close()

	// Return the build id because the response redirects
	u := res.Request.URL
	segments := strings.Split(u.Path, "/")
	return segments[len(segments)-1]
}

func (s *SolderClient) GetBuildId(modpack *types.Modpack) string {
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

func (s *SolderClient) GetModpackId(slug string) string {
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

func (s *SolderClient) GetModId(modid string) string {
	if id, ok := s.modIdCache[modid]; ok {
		if id != "" {
			return id
		}
	}

	Url := s.createUrl("mod/list")

	response := s.doRequest(http.MethodGet, Url.String(), "")
	defer response.Body.Close()
	mods := crawlers.CrawlModList(response)
	//log.Panic(mods)
	for _, mod := range mods {
		if strings.ToLower(mod.Name) == strings.ToLower(modid) {
			id := mod.Id
			if strings.Trim(id, " ") != "" {
				s.modIdCache[modid] = id
			}
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

func (s *SolderClient) AddModversionToBuild(mod *types.OutputInfo, modpackBuildId string) {
	Url := s.createUrl("modpack/modify/add")

	form := url.Values{}
	form.Add("build", modpackBuildId)
	form.Add("mod-name", mod.Id)
	form.Add("mod-version", mod.GenerateOnlineVersion())
	form.Add("action", "add")

	response := s.postForm(Url.String(), form)
	defer response.Body.Close()

	var res addBuildToModpackResponse
	err := json.NewDecoder(response.Body).Decode(&res)

	if err != nil {
		log.Println(Url.String())
		log.Printf("%v", form)
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		log.Println(buf.String())
		log.Panic(err)
	}

	if res.Status != "success" {
		log.Println(res.Reason)
		log.Println(*mod)
		log.Panic("Something went wrong when adding a mod to a build, see above mod details")
	}

}
