package coll33tx

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"unicode/utf8"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type Item struct {
	Name    string
	Href    *url.URL
	Size    string
	Seeds   int
	Leeches int
}

type ListPageCrawledCallback func(items []*Item, currentPage int, nextPage *url.URL, response *colly.Response, log *log.Entry)

// ListCollector is a wrapper around the colly collector + listing page data
type ListCollector struct {
	*colly.Collector
	Data          map[int][]*Item // page => items
	TotalPages    int
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
		fmt.Printf("%s", string(r.Body))
		/*
			todo uncomment these when tests for the 2 funcs below pass
			doc, err := getResponseDoc(r)
			if err != nil {
				col.Log.WithError(err).Error("couldn't parse page document")
				return
			}

			items, errItems := getPageItems(doc, r.Request)
			if errItems != nil {
				col.Log.WithError(err).Error("didn't find any list elements")
			}
			pagesCount, errPagesCount := getPagination(doc)
			if errPagesCount != nil {
				col.Log.WithError(errPagesCount).Warn("missing pagesCount")
			}
			currentPage, errCurrentPage := getCurrentPage(doc)
			if errCurrentPage != nil {
				col.Log.WithError(errCurrentPage).Warn("missing currentPage")
			}

			col.OnPageCrawled(items, currentPage, pagesCount, r, col.Log)
		*/
	})

	return &col
}

type listPageDocument struct {
	*goquery.Document
}

func NewListPageDocument(r *colly.Response) (*listPageDocument, error) {
	d, err := getResponseDoc(r)
	if err != nil {
		return nil, err
	}

	return &listPageDocument{d}, nil
}

func (doc *listPageDocument) GetPagination() (pagination struct {
	totalPages, currentPage int
}, err error) {
	selector := ".pagination"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("no %s", selector)
		return
	}

	getPageNumberFromLink := func(link *goquery.Selection) (page int, err error) {
		if link.Nodes == nil {
			return page, fmt.Errorf("missing pagination link")
		}

		href, exists := link.Attr("href")
		if !exists {
			return page, fmt.Errorf("missing href attr for pagination link")
		}

		u, err := doc.Url.Parse(href)
		if err != nil {
			return page, fmt.Errorf("couldn't parse href for pagination link")
		}

		pieces := strings.Split(u.String(), "/")

		if len(pieces) < 2 {
			return page, fmt.Errorf("can't extract page number from url %s in pagination link", u.String())
		}

		return strconv.Atoi(pieces[len(pieces)-2])
	}

	pagination.totalPages, err = getPageNumberFromLink(selection.Find(".last a"))
	if err != nil {
		return
	}

	pagination.currentPage, err = getPageNumberFromLink(selection.Find(".active a"))
	if err != nil {
		return
	}

	return
}

// GetPageItems returns list items along with their errors / missing stuff
func (doc *listPageDocument) GetPageItems() (items []struct {
	item *Item
	errs []error
}, err error) {
	trs := doc.Find(".page-content .table-striped tbody tr")
	if trs.Nodes == nil {
		err = errors.New("couldn't find list trs")
		return
	}

	trToItem := func(i int, tr *goquery.Selection) (item *Item, errs []error) {
		item = new(Item)
		var err error

		td := tr.Find("td.name")
		if td.Nodes != nil {
			item.Name = td.Text()
		}
		aName := td.Find("a").Eq(1)
		if aName.Nodes != nil {
			href, exists := aName.Attr("href")
			if !exists {
				errs = append(errs, fmt.Errorf("list item %d (%s) has no link behind", i, item.Name))
			}
			item.Href, err = doc.Url.Parse(href)
			if err != nil {
				errs = append(errs, fmt.Errorf("link behind link %d (%s) is invalid", i, item.Name))
			}
		}

		td = tr.Find("td.seeds")
		if td.Nodes != nil {
			if item.Seeds, err = strconv.Atoi(td.Text()); err != nil {
				errs = append(errs, fmt.Errorf("can't convert seeders to int for item %d (%s)", i, item.Name))
			}
		}

		td = tr.Find("td.leeches")
		if td.Nodes != nil {
			if item.Leeches, err = strconv.Atoi(td.Text()); err != nil {
				errs = append(errs, fmt.Errorf("can't convert leechers to int for item %d (%s)", i, item.Name))
			}
		}

		td = tr.Find("td.size")
		if td.Nodes != nil {
			unwantedSpan := td.Find("span")
			if unwantedSpan.Nodes != nil {
				unwantedSpan.Remove()
			}
			item.Size = td.Text()
			if utf8.RuneCountInString(item.Size) < 3 {
				errs = append(errs, fmt.Errorf("weird torrent size for item %d (%s)", i, item.Name))
			}
		}

		return item, errs
	}

	trs.Each(func(i int, tr *goquery.Selection) {
		item, errs := trToItem(i, tr)
		items = append(items, struct {
			item *Item
			errs []error
		}{item: item, errs: errs})
	})

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
