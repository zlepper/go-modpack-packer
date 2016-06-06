package crawlers

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type ModVersion struct {
	Filesize string
	Md5      string
	Url      string
	Id       string
	Version  string
}

func CrawlModVersion(res *http.Response) []ModVersion {
	mvs := make([]ModVersion, 0)

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Panic(err)
	}

	doc.Find("table > tbody > tr.version").Each(func(_ int, sel *goquery.Selection) {
		var mv ModVersion
		mv.Id, _ = sel.Attr("rel")
		mv.Version = sel.Find("td.version").Text()
		mv.Url = sel.Find("td.url").Text()
		mv.Md5 = sel.Find("td > input").AttrOr("placeholder", "")
		mvs = append(mvs, mv)
	})

	return mvs
}
