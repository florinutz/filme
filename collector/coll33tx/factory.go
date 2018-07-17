package coll33tx

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

func (col *ListCollector) CanHandleResponse(r *colly.Response) bool {
	doc, err := getDoc(r)
	if err != nil {
		log.WithField("response", r).Warn("error while checking if ListCollector can handle a response")
		return false
	}
	if listElement := doc.Find(".page-content .featured-list"); listElement.Nodes == nil {
		return false
	}
	return true
}

func (col *DetailsCollector) CanHandleResponse(r *colly.Response) bool {
	doc, err := getDoc(r)
	if err != nil {
		log.WithField("response", r).Warn("error while checking if ListCollector can handle a response")
		return false
	}
	if title := doc.Find(".box-info-heading h1"); title.Nodes == nil {
		return false
	}
	if img := doc.Find(".torrent-detail .torrent-image img"); img.Nodes == nil {
		return false
	}
	return true
}

func getDoc(r *colly.Response) (doc *goquery.Document, err error) {
	doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	return
}
