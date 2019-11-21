package list

import (
	"fmt"
	"net/url"
	"strconv"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/florinutz/filme/pkg/collector"
	"github.com/gocolly/colly"
)

const (
	TestPageList               = "https://1337x.to/popular-movies"
	TestPageListWithPagination = "https://1337x.to/search/romania/3/"
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

// GetLines returns list items along with their errors / missing stuff
func (doc *document) GetLines() ([]*Line, error) {
	selector := ".page-content .table-striped tbody tr"
	trs := doc.Find(selector)
	if trs.Nodes == nil {
		return nil, fmt.Errorf("selector '%s' not found in document at url %s", selector, doc.Url.String())
	}

	var lines []*Line
	trs.Each(func(i int, tr *goquery.Selection) {
		line := new(Line)
		line.Item, line.Errs = doc.trToItem(i, tr)
		lines = append(lines, line)
	})

	return lines, nil
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
