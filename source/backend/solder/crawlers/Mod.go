package crawlers

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Mod struct {
	Id          string
	Name        string
	PrettyName  string
	Author      string
	Description string
	Link        string
	Donate      string
	Versions    []string
	Active      string
}

const complexNamePattern string = `(.+?) ?\((.+?)\)(?:.+)`

func CrawlModList(res *http.Response) []Mod {
	doc := makeDoc(res)
	mods := make([]Mod, 0)

	tableRows := doc.Find("table > tbody > tr")

	c := make(chan Mod)
	tableRows.Each(func(_ int, r *goquery.Selection) {
		go func(row *goquery.Selection) {

			var mod Mod

			mod.Id = row.Find(" td:nth-child(1)").Text()

			// Read name and slug
			content := row.Find(" td:nth-child(2)").Text()
			// Remove newlines
			re := regexp.MustCompile("\\r\\n?|\\n|\\t")
			content = string(re.ReplaceAll([]byte(content), []byte("")))

			// Remove double spaces
			content = strings.Replace(content, "  ", " ", -1)
			// Get matches
			re = regexp.MustCompile(complexNamePattern)
			r := re.FindStringSubmatch(content)
			if len(r) > 2 {
				mod.PrettyName = r[1]
				mod.Name = r[2]
			} else {
				log.Println(content)
				log.Println(r)
				log.Panic("Something went wrong when regexing stuff")
			}

			c <- mod

		}(r)
	})
	for i := 0; i < tableRows.Length(); i++ {
		mod := <-c
		mods = append(mods, mod)
	}

	return mods
}

func CrawlMod(res *http.Response) Mod {
	doc := makeDoc(res)

	mod := Mod{}

	mod.PrettyName = doc.Find("#pretty_name").AttrOr("value", "")
	mod.Name = doc.Find("#name").AttrOr("value", "")
	mod.Author = doc.Find("#author").AttrOr("value", "")
	mod.Description = doc.Find("#description").Text()
	mod.Link = doc.Find("#link").Text()
	mod.Donate = doc.Find("#donatelink").AttrOr("value", "")

	mod.Id = doc.Find(`#add-row > input[name="mod-id"]`).AttrOr("value", "")
	doc.Find("#versions > table > tbody > tr.version > td.version").Each(func(index int, s *goquery.Selection) {
		mod.Versions = append(mod.Versions, s.Text())
	})

	return mod
}
