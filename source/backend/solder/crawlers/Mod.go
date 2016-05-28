package crawlers

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"go/build"
	"golang.org/x/net/html"
	"regexp"
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


const complexNamePattern string = "(?:<.+?>)(.+?)(?:<.+?>) ?\\((.+?)\\)(?:.+)"

func CrawlModList(res *http.Response) []Mod {
	doc := makeDoc(res)

	mods := make([]Mod, 0)

	tableRows := doc.Find("table > tbody > tr")

	c := make(chan Mod)

	for _, node := range tableRows.Nodes {
		go func(n html.Node) {

			var mod Mod
			row := newSingleSelection(n, doc)

			mod.Id = row.Find(" td:nth-child(1)").Text()

			// Read name and slug
			content := row.Find(" td:nth-child(2)").Text()
			// Remove newlines
			re := regexp.MustCompile("\\r\\n?|\\n|\\t")
			re.ReplaceAll(content, []byte(""))

			// Remove double spaces
			content = strings.Replace(content, "  ", " ", -1)

			// Get matches
			re = regexp.MustCompile(complexNamePattern)
			r := re.FindAllString(content, -1)
			mod.PrettyName = r[0]
			mod.Name = r[1]

			c <- mod

		}(node)
	}

	for i := 0; i < tableRows.Length(); i++ {
		mod := <- c
		mods = append(mods, mod)
	}

	return mods
}
