package solder

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"log"
	"path"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/crawlers"
	"github.com/zlepper/go-modpack-packer/source/backend/handlers"
	"strings"
	"encoding/json"
)


type solderClient struct {
	Client http.Client
	Url url.URL
	modVersionIdCache map[string]map[string]string
	modIdCache map[string]string
	buildCache map[string]crawlers.Build
}

type SolderClient *solderClient

func NewSolderClient(Url string) SolderClient {
	var cookieJar, _ = cookiejar.New(nil)

	var client = &http.Client{
		Jar: cookieJar,
	}

	u, err := url.Parse(Url)
	if err != nil {
		log.Panic(err)
	}


	return &solderClient{
		Client: client,
		Url: u,
	}
}

func (s *solderClient) createUrl(after string) url.URL {
	var url url.URL
	*url = *s.Url
	url.Path = path.Join(url.Path, after)
	return url
}

func (s *solderClient) doRequest(method, url, data string) (*http.Response) {
	req, _ := http.NewRequest(method, url, data)
	response, err := s.Client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	return response
}

func (s *solderClient) Login(email string, password string) (bool) {
	Url := s.createUrl("login")

	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)

	response := s.doRequest(http.MethodPost, Url.String(), form.Encode())
	defer response.Close()

	return crawlers.CrawlLogin(response)
}

func (s *solderClient) CreatePack(name, slug string) string {

}

func (s *solderClient) AddMod(mod handlers.Mod) string {
	Url := s.createUrl("mod/add-version")

	form := url.Values{}
	form.Add("pretty_name", mod.Name)
	form.Add("name", mod.ModId)
	form.Add("author", mod.Authors)
	form.Add("description", mod.Description)
	form.Add("link", mod.Url)
	response := s.doRequest(http.MethodPost, Url.String(), form.Encode())
	defer response.Close()

	// TODO Return
}

func (s *solderClient) GetActiveModversionInBuildId(mod handlers.Mod, buildId string) string {
	Url := s.createUrl("modpack/build/" + buildId)
	res := s.doRequest(http.MethodGet, Url.String(), "")
	defer res.Close()
	build := crawlers.CrawlBuild(res)

}

func (s *solderClient) SetModVersionInBuild(mod handlers.Mod, buildId string) {
	modVersionId := s.GetModVersionId(mod)
	Url := s.createUrl("modpack/build/modify")

	form := url.Values{}
	form.Add("action", "version")
	form.Add("build-id", buildId)
	form.Add("version", modVersionId)
	form.Add("modversion-id", s.GetActiveModversionInBuildId(mod, buildId))

	res := s.doRequest(http.MethodPost, Url.String(), form.Encode())
	defer res.Close()

}

func (s *solderClient) IsPackOnline(modpack handlers.Modpack) bool {
	return s.GetModpackId(modpack.GetSlug()) != ""
}

func (s *solderClient) IsBuildOnline(modpack handlers.Modpack) bool {
	return s.GetBuildId(modpack) != ""
}

func (s *solderClient) GetModVersionId(mod handlers.Mod) string {
	modId, modVersion := mod.ModId, mod.GenerateOnlineVersion()
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

func (s *solderClient) CreateBuild(modpack handlers.Modpack) string {
	Url := s.createUrl("modpack/add-build/" + s.GetModpackId(modpack.GetSlug()))

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
	return segments[len(segments) - 1]
}

func (s *solderClient) GetBuildId(modpack handlers.Modpack) string {
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

func (s *solderClient) GetModpackId(slug string) string {
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


func (s *solderClient) GetModId(modid string) string {
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

func (s *solderClient) AddModversionToBuild(mod handlers.Mod, modpackBuildId string) {
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
