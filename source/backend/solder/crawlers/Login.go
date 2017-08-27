package crawlers

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

func CrawlLogin(response *http.Response) error {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return err
	}

	// Find the dashboard <h1> node.
	// This node is only available if login was successful
	nodes := doc.Find("#page-wrapper > div > div > div:first-child > h1").Nodes

	// If some nodes was found, then login was successful
	if len(nodes) > 0 {
		return nil
	} else {
		return errors.New("Login attempt failed")
	}
}
