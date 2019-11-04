package collector

import (
	"bytes"
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// GetResponseDoc converts a colly response to a goquery document
func GetResponseDocument(r *colly.Response) (doc *goquery.Document, err error) {
	document, err := goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	if document == nil {
		return nil, errors.New("this should not be nil")
	}
	document.Url = r.Request.URL

	return document, err
}
