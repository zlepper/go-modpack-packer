package handlers

import (
	"encoding/json"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/crawlers"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type SolderClient struct {
	Client            http.Client
	Url               url.URL
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
		Client: client,
		Url:    u,
	}
}

func (s *SolderClient) createUrl(after string) url.URL {
	var url url.URL
	*url = *s.Url
	url.Path = path.Join(url.Path, after)
	return url
}

func (s *SolderClient) doRequest(method, url, data string) *http.Response {
	req, _ := http.NewRequest(method, url, data)
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

	response := s.doRequest(http.MethodPost, Url.String(), form.Encode())
	defer response.Close()

	return crawlers.CrawlLogin(response)
}

func (s *SolderClient) CreatePack(name, slug string) string {
	id := s.GetModpackId(slug)
	if id != "" {
		return id
	}
	form := url.Values{}
	form.Add("name", name)
	form.Add("slug", slug)
	form.Add("hidden", false)
	res := s.doRequest(http.MethodPost, "modpack/create", form.Encode())
	defer res.Close()

	u, _ := res.Location()
	segments := strings.Split(u.Path, "/")
	segment := segments[len(segments)-1]
	if _, err := strconv.Atoi(segment); err == nil {
		return segment
	} else {
		return s.GetModpackId(slug)
	}

}

func (s *SolderClient) AddMod(mod *outputInfo) string {
	Url := s.createUrl("mod/add-version")

	form := url.Values{}
	form.Add("pretty_name", mod.Name)
	form.Add("name", mod.Id)
	form.Add("author", mod.Author)
	form.Add("description", mod.Description)
	form.Add("link", mod.Url)
	response := s.doRequest(http.MethodPost, Url.String(), form.Encode())
	defer response.Close()

	return s.GetModId(mod.Id)
}

func (s *SolderClient) AddModVersion(modId, md5, version string) {
	form := url.Values{}
	form.Add("mod-id", modId)
	form.Add("add-version", version)
	form.Add("add-md5", md5)
	res := s.doRequest(http.MethodPost, "mod/add-version", form.Encode())
	res.Close()
}

func (s *SolderClient) RehashModVersion(modversionId string, md5 string) {
	form := url.Values{}
	form.Add("version-id", modversionId)
	form.Add("md5", md5)
	res := s.doRequest(http.MethodPost, "mod/rehash", form.Encode())
	res.Close()
}

func (s *SolderClient) IsModversionOnline(mod *outputInfo) bool {
	return s.GetModVersionId(mod) != ""
}

func (s *SolderClient) GetBuild(buildId string) crawlers.Build {
	if build, ok := s.buildCache[buildId]; ok {
		return build
	}

	Url := s.createUrl("modpack/build/" + buildId)
	res := s.doRequest(http.MethodGet, Url.String(), "")
	defer res.Close()
	build := crawlers.CrawlBuild(res)

	s.buildCache[buildId] = build
	return build
}

func (s *SolderClient) IsModInBuild(mod *outputInfo, buildId string) bool {
	build := s.GetBuild(buildId)

	for _, m := range build.Mods {
		if m.Name == mod.Id {
			return true
		}
	}
	return false
}

func (s *SolderClient) IsModversionActiveInBuild(mod *outputInfo, buildId string) bool {
	build := s.GetBuild(buildId)

	for _, m := range build.Mods {
		if m.Name == mod.Id && m.Active == mod.GenerateOnlineVersion() {
			return true
		}
	}
	return false
}

func (s *SolderClient) GetActiveModversionInBuildId(mod *outputInfo, buildId string) string {
	Url := s.createUrl("modpack/build/" + buildId)
	res := s.doRequest(http.MethodGet, Url.String(), "")
	defer res.Close()
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

func (s *SolderClient) SetModVersionInBuild(mod *outputInfo, buildId string) {
	modVersionId := s.GetModVersionId(mod)
	Url := s.createUrl("modpack/build/modify")

	form := url.Values{}
	form.Add("action", "version")
	form.Add("build-id", buildId)
	form.Add("version", modVersionId)
	form.Add("modversion-id", s.GetActiveModversionInBuildId(mod, buildId))

	res := s.doRequest(http.MethodPost, Url.String(), form.Encode())
	res.Close()
}

func (s *SolderClient) IsPackOnline(modpack *Modpack) bool {
	return s.GetModpackId(modpack.GetSlug()) != ""
}

func (s *SolderClient) IsBuildOnline(modpack *Modpack) bool {
	return s.GetBuildId(modpack) != ""
}

func (s *SolderClient) GetModVersionId(mod *outputInfo) string {
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
	defer res.Close()
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

func (s *SolderClient) CreateBuild(modpack *Modpack, modpackId string) string {
	Url := s.createUrl("modpack/add-build/" + modpackId)

	form := url.Values{}
	form.Add("version", modpack.Version)
	form.Add("minecraft", modpack.MinecraftVersion)
	form.Add("java-version", modpack.Java)
	form.Add("memory-enabled", modpack.Memory != "")
	form.Add("momory", modpack.Memory)
	res := s.doRequest(http.MethodPost, Url.String(), form.Encode())
	defer res.Close()

	// Return the build id because the response redirects (TODO Test if go also follow redirect
	u, _ := res.Location()
	segments := strings.Split(u.Path, "/")
	return segments[len(segments)-1]
}

func (s *SolderClient) GetBuildId(modpack *Modpack) string {
	Url := s.createUrl("modpack/view/" + modpack.GetSlug())

	res := s.doRequest(http.MethodGet, Url.String(), "")
	builds := crawlers.CrawlBuildList(res)

	for _, build := range builds {
		if build.Version == modpack.Version {
			return build.Id
		}
	}
	return ""
}

func (s *SolderClient) GetModpackId(slug string) string {
	Url := s.createUrl("modpack/list")

	res := s.doRequest(http.MethodGet, Url.String(), "")
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
		return id
	}

	Url := s.createUrl("mod/list")

	response := s.doRequest(http.MethodGet, Url.String(), "")
	defer response.Close()
	mods := crawlers.CrawlModList(response)
	for _, mod := range mods {
		if mod.Name == modid {
			id := mod.Id
			if strings.Trim(id, " ") != "" {
				s.modIdCache[modid] = id
			}
			return id
		}
	}
	return ""
}

type addBuildToModpackResponse struct {
	Status string `json:"status"`
}

func (s *SolderClient) AddModversionToBuild(mod *outputInfo, modpackBuildId string) {
	Url := s.createUrl("modpack/modify/add")

	form := url.Values{}
	form.Add("build", modpackBuildId)
	form.Add("mod-name", mod.Name)
	form.Add("mod-version", mod.GenerateOnlineVersion())
	form.Add("action", "add")

	response := s.doRequest(http.MethodPost, Url.String(), form.Encode())
	defer response.Close()

	var res addBuildToModpackResponse
	err := json.NewDecoder(response.Body).Decode(&res)

	if err != nil {
		log.Panic(err)
	}

	if res["status"].(string) != "success" {
		log.Panic("Something went wrong when adding a mod to a build")
	}

}
