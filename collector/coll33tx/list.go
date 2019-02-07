package coll33tx

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"

	"strconv"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type Item struct {
	Name    string
	Href    string
	Size    string
	Seeds   int
	Leeches int
}

type ListPageCrawledCallback func(items []*Item, currentPage int, nextPage *url.URL, response *colly.Response, log *log.Entry)

// ListCollector is a wrapper around the colly collector + listing page data
type ListCollector struct {
	*colly.Collector
	Data          map[int][]*Item // page => items
	totalPages    int
	OnPageCrawled ListPageCrawledCallback
	Log           *log.Entry
}

func NewListCollector(
	onPageCrawled ListPageCrawledCallback,
	log *log.Entry,
	options ...func(collector *colly.Collector),
) *ListCollector {
	col := ListCollector{
		Collector:     initCollector(log, options...),
		OnPageCrawled: onPageCrawled,
		Log:           log,
	}
	col.Collector.OnScraped(func(r *colly.Response) {
		doc, err := getResponseDoc(r)
		if err != nil {
			col.Log.WithError(err).Error("couldn't parse page document")
			return
		}
		items, errItems := getPageItems(doc, r.Request)
		if errItems != nil {
			col.Log.WithError(err).Error("didn't find any list elements")
		}
		nextPage, errNextPage := getNextPage(doc)
		if errNextPage != nil {
			col.Log.WithError(errNextPage).Warn("missing nextPage")
		}
		currentPage, errCurrentPage := getCurrentPage(doc)
		if errCurrentPage != nil {
			col.Log.WithError(errCurrentPage).Warn("missing currentPage")
		}

		col.OnPageCrawled(items, currentPage, nextPage, r, col.Log)
	})

	return &col
}

func getCurrentPage(doc *goquery.Document) (page int, err error) {
	selector := ".pagination li.active"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("couldn't find the current page: no element at selector '%s'", selector)
		return
	}
	return strconv.Atoi(selection.Text())
}

func getNextPage(doc *goquery.Document) (url *url.URL, err error) {
	selector := ".pagination li.active"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("couldn't find the current page: no element at selector '%s'", selector)
		return
	}
	return
}

func getPageItems(doc *goquery.Document, req *colly.Request) (items []*Item, err error) {
	selector := "td.name a:nth-of-type(2)"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("couldn't select page items using selector '%s'", selector)
		return
	}
	for _, node := range selection.Nodes {
		item := Item{
			Name: node.Data,
			// todo href?
		}
		items = append(items, &item)
	}
	return
}

/*func (i *Item) fromTitleLink(e *colly.HTMLElement) (errs []error) {
	i.Name = e.Text
	i.Href = e.Request.AbsoluteURL(e.Attr("href"))

	tr := e.DOM.Parent().Parent()

	if seeds, err := strconv.Atoi(tr.Find(".seeds").Text()); err != nil {
		errs = append(errs, err)
	} else {
		i.Seeds = seeds
	}

	if leeches, err := strconv.Atoi(tr.Find(".leeches").Text()); err != nil {
		errs = append(errs, err)
	} else {
		i.Leeches = leeches
	}

	i.Size = tr.Find(".size").Text()

	return
}

func (col *ListCollector) CanHandleResponse(r *colly.Response) bool {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	if err != nil {
		log.WithField("response", r).Warn("error while checking if ListCollector can handle a response")
		return false
	}

	if listElement := doc.Find(".page-content .featured-list"); listElement.Nodes == nil {
		return false
	}

	return true
}*/
