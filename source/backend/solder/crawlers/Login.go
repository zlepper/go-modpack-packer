package crawlers

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"log"
)

func CrawlLogin(response *http.Response) bool {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		log.Panic(err)
	}

	// Find the dashboard <h1> node.
	// This node is only available if login was successful
	nodes := doc.Find("#page-wrapper > div > div > div:first-child > h1").Nodes

	// If some nodes was found, then login was successful
	return len(nodes) > 0
}
