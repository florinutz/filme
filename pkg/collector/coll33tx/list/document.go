package list

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/florinutz/filme/pkg/collector"
	"github.com/gocolly/colly"
)

// document describes a list page
type document struct {
	*goquery.Document
}

func NewDocument(r *colly.Response) (*document, error) {
	d, err := collector.GetResponseDocument(r)
	if err != nil {
		return nil, err
	}
	return &document{d}, nil
}

type Item struct {
	Name    string
	Href    *url.URL
	Size    string
	Seeders int
	Leeches int
}

type Line struct {
	Item *Item
	Errs []error
}

type Pagination struct {
	PagesCount  int
	CurrentPage int
	Err         error
}

// GetPagination returns nil if pagination is not present in the page and 0 for both Pagination values
// if they're not found
func (doc *document) GetPagination() *Pagination {
	selector := ".pagination"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		return nil
	}
	pagination := new(Pagination)
	if page, err := doc.getPageNumberFromLink(selection.Find(".last a")); err == nil {
		pagination.PagesCount = page
	}
	if page, err := doc.getPageNumberFromLink(selection.Find(".active a")); err == nil {
		pagination.CurrentPage = page
	}
	return pagination
}

func (doc *document) getPageNumberFromLink(link *goquery.Selection) (page int, err error) {
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

// GetLines returns list items along with their errors / missing stuff
func (doc *document) GetLines() (lines []*Line) {
	trs := doc.Find(".page-content .table-striped tbody tr")
	if trs.Nodes == nil {
		return
	}
	trs.Each(func(i int, tr *goquery.Selection) {
		line := new(Line)
		line.Item, line.Errs = doc.trToItem(i, tr)
		lines = append(lines, line)
	})
	return
}

func (doc *document) trToItem(i int, tr *goquery.Selection) (item *Item, errs []error) {
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
		if item.Seeders, err = strconv.Atoi(td.Text()); err != nil {
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
