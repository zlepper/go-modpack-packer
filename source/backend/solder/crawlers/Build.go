package crawlers

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"log"
	"golang.org/x/net/html"
	"strings"
)

type Build struct {
	Id string
	Minecraft string
	Java string
	Memory string
	Mods []Mod
	Version string
}

func CrawlBuild(res *http.Response) (Build) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Panic(err)
	}

	var build Build

	javaVersion := doc.Find("#page-wrapper > div > div > div:nth-child(2) > div:nth-child(2) > div:nth-child(2) > label:first-child > span").Text()
	if javaVersion != "Not Required" {
		build.Java = javaVersion
	}

	memory := doc.Find("#page-wrapper > div > div > div:nth-child(2) > div:nth-child(2) > div:nth-child(2) > label:nth-child(2) > span").Text()
	if memory != "Not Required" {
		build.Memory = memory
	}

	build.Mods = make([]Mod, 0)

	// Find all the mods in the build
	tableRows := doc.Find("table#mod-list > tbody > tr")

	c := make(chan Mod)

	for _, node := range tableRows.Nodes {
		go func(n html.Node) {

			var mod Mod
			row := newSingleSelection(n, doc)

			// Use regex to calculate modname and slug
			firstPart := row.Find("td:first-child").First().Text()
			matches := re.FindAllString(firstPart, -1)
			mod.PrettyName = matches[0]
			mod.Name = matches[1]

			// Find mod id
			anchor := row.Find("a")
			url, _ := anchor.Attr("href")
			mod.Id = url[strings.LastIndex(url, "/")+ 1:]

			// Find mod versions
			mod.Versions = make([]string, 0)
			row.Find("select > option").Each(func(_, s *goquery.Selection) {
				_, exists := s.Attr("selected")
				if exists {
					mod.Active = s.Text()
				}
				mod.Versions = append(mod.Versions, s.Text())
			})

			c <- mod

		}(node)
	}

	for i := 0; i < tableRows.Length(); i++ {
		mod := <- c
		build.Mods = append(build.Mods, mod)
	}

	return build
}

func CrawlBuildList(res *http.Response) ([]Build) {

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Panic(err)
	}

	builds := make([]Build, 0)

	tableRows := doc.Find("table#dataTables > tbody > tr")

	tableRows.Each(func(_, row *goquery.Selection) {
		build := Build{
			Id: row.Find("td:nth-child(1)").Text(),
			Minecraft: row.Find("td:nth-child(2)").Text(),
			Version: row.Find("td:nth-child(3)").Text(),
		}
		builds = append(builds, build)
	})

	return builds
}

// Helper constructor to create a selection of only one node
func newSingleSelection(node *html.Node, doc *goquery.Document) *goquery.Selection {
	return &goquery.Selection{[]*html.Node{node}, doc, nil}
}
