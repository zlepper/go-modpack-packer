package crawlers

import (
	"regexp"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"log"
)

var re *regexp.Regexp

func init() {
	const namePattern string = "([^\\r\\n\\(\\)]+?) \\(([^ \\r\\n\\(\\)]+)\\)"
	re = regexp.MustCompile(namePattern)
}

func makeDoc(res *http.Response) (goquery.Document) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Panic(err)
	}
	return doc
}
