package crawlers

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type Modpack struct {
	Id          string
	Name        string
	DisplayName string
	Recommended string
	Latest      string
	Build       []Build
}

func CrawlModpackList(res *http.Response) []Modpack {
	doc := makeDoc(res)

	modpacks := make([]Modpack, 0)

	tableRows := doc.Find("table > tbody >tr")

	tableRows.Each(func(_ int, row *goquery.Selection) {
		var modpack Modpack
		modpack.DisplayName = row.Find("td:nth-child(1)").Text()
		modpack.Name = row.Find("td:nth-child(2)").Text()
		modpack.Recommended = row.Find("td:nth-child(3)").Text()
		modpack.Latest = row.Find("td:nth-child(4)").Text()

		tid, _ := row.Find("td:nth-child(7) > a:first-child").Attr("href")
		tid = tid[strings.LastIndex(tid, "/")+1:]
		modpack.Id = tid
		modpacks = append(modpacks, modpack)
	})

	return modpacks
}
