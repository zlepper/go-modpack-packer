package crawlers

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"log"
)

type ModVersion struct {
	Filesize string
	Md5 string
	Url string
	Id string
	Version string
}


func CrawlModVersion(res *http.Response) []ModVersion {
	mvs := make([]ModVersion, 0)

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Panic(err)
	}

	doc.Find("table > tbody > tr.version").Each(func(_, sel *goquery.Selection) {
		var mv ModVersion
		mv.Id, _ = sel.Attr("rel")
		mv.Version = sel.Find("td.version").Text()
		mv.Url = sel.Find("td.url").Text()
		mv.Md5 = sel.Find("td > input").AttrOr("placeholder", "")
		mvs = append(mvs, mv)
	})

	return mvs
}

