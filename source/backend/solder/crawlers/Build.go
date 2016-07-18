package crawlers

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

type Build struct {
	Id        string
	Minecraft string
	Java      string
	Memory    string
	Mods      []Mod
	Version   string
}

func CrawlBuild(res *http.Response) Build {
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
	tableRows.Each(func(_ int, r *goquery.Selection) {
		go func(row *goquery.Selection) {

			var mod Mod

			// Use regex to calculate modname and slug
			firstPart := row.Find("td:first-child").First().Text()
			matches := re.FindStringSubmatch(firstPart)
			log.Println(matches)
			mod.PrettyName = matches[1]
			mod.Name = matches[2]

			// Find mod id
			anchor := row.Find("a")
			url, _ := anchor.Attr("href")
			mod.Id = url[strings.LastIndex(url, "/")+1:]

			// Find mod versions
			mod.Versions = make([]string, 0)
			row.Find("select > option").Each(func(_ int, s *goquery.Selection) {
				_, exists := s.Attr("selected")
				if exists {
					mod.Active = s.Text()
				}
				mod.Versions = append(mod.Versions, s.Text())
			})

			c <- mod
		}(r)
	})
	for i := 0; i < tableRows.Length(); i++ {
		mod := <-c
		build.Mods = append(build.Mods, mod)
	}

	return build
}

func CrawlBuildList(res *http.Response) []Build {

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Panic(err)
	}

	builds := make([]Build, 0)

	tableRows := doc.Find("table#dataTables > tbody > tr")

	tableRows.Each(func(_ int, row *goquery.Selection) {
		build := Build{
			Id:        row.Find("td:nth-child(1)").Text(),
			Minecraft: row.Find("td:nth-child(3)").Text(),
			Version:   row.Find("td:nth-child(2)").Text(),
		}
		builds = append(builds, build)
	})

	return builds
}
