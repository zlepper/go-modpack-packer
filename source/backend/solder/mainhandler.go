package solder

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"log"
	"path"
	"github.com/zlepper/go-modpack-packer/source/backend/solder/crawlers"
	"github.com/zlepper/go-modpack-packer/source/backend/handlers"
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

func (s *solderClient) GetModId(modid string) string {
	if id, ok := s.modIdCache[modid]; ok {
		return id
	}

	Url := s.createUrl("mod/list")

	response := s.doRequest(http.MethodGet)
}
