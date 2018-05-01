package solder

import (
	"net/http"
	"net/url"
	"path"
	"log"
	"encoding/json"
	"io"
	"bytes"
)

type SolderClientV8 struct {
	client *http.Client
	url    *url.URL
	apiKey string
}

type PackageItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type packageListResponse struct {
	Data []PackageItem `json:"data"`
}

func (c *SolderClientV8) createUrl(after string) url.URL {
	url := *c.url
	url.Path = path.Join(url.Path, after)
	log.Println(url.Path)
	return url
}

func (c *SolderClientV8) createGetRequest(url url.URL) (*http.Request, error) {
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+c.apiKey)
	request.Header.Add("Accept", "application/json")

	return request, nil
}

func (c *SolderClientV8) createPostRequest(url url.URL, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest("POST", url.String(), body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+c.apiKey)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

// Test the solder connection
func (c *SolderClientV8) Test() error {
	_, err := c.GetPackages()
	return err
}

func (c *SolderClientV8) GetPackages() ([]PackageItem, error) {
	packagesUrl := c.createUrl("/api/packages")

	request, err := c.createGetRequest(packagesUrl)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data := new(packageListResponse)
	err = json.NewDecoder(response.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}

type PackageCreateRequest struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	ProjectUrl  string `json:"project_url"`
	DonationUrl string `json:"donation_url"`
}

type packageCreateResponse struct {
	Data PackageItem `json:"data"`
}

func (c *SolderClientV8) CreatePackage(args PackageCreateRequest) (*PackageItem, error) {
	body, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	packagesUrl := c.createUrl("/api/packages")

	request, err := c.createPostRequest(packagesUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data := new(packageCreateResponse)
	err = json.NewDecoder(response.Body).Decode(data)
	if err != nil {
		return nil, err
	}

	return &data.Data, nil
}

//func (c *SolderClientV8) GetPackageVersions(packageId string) ()
