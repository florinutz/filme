package collector

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// GetResponseDoc converts a colly response to a goquery document
func GetResponseDocument(r *colly.Response) (doc *goquery.Document, err error) {
	document, err := goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	document.Url = r.Request.URL

	return document, err
}
