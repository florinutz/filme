package list

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const ItemsPerPage = 20

var urlTplSplitter = regexp.MustCompile(`(?m)/\d/$`)

type Pagination struct {
	PagesCount int
	Current    int
	CurrentUrl *url.URL
	NextUrl    *url.URL
	LinksTpl   string
}

// GetPagination returns nil if pagination is not present in the page and 0 for both Pagination values
// if they're not found
func (doc *document) GetPagination() *Pagination {
	selection := doc.Find(".pagination")
	if selection.Nodes == nil {
		return nil
	}
	pagination := &Pagination{
		CurrentUrl: doc.Url,
	}

	pagination.PagesCount, _, _ = doc.readPageLink(selection.Find(".last a"))
	pagination.Current, pagination.NextUrl, _ = doc.readPageLink(selection.Find(".active a"))

	if pagination.NextUrl != nil {
		pagination.LinksTpl = urlTplSplitter.ReplaceAllString(pagination.NextUrl.String(), "/%d/")
	}

	return pagination
}

func (doc *document) readPageLink(link *goquery.Selection) (page int, nextPageUrl *url.URL, err error) {
	if link.Nodes == nil {
		return page, nil, fmt.Errorf("missing pagination link")
	}

	href, exists := link.Attr("href")
	if !exists {
		return page, nil, fmt.Errorf("missing href attr for pagination link")
	}

	pageUrl, err := doc.Url.Parse(href)
	if err != nil {
		return page, nil, fmt.Errorf("couldn't parse href for pagination link")
	}

	pieces := strings.Split(pageUrl.String(), "/")

	if len(pieces) < 2 {
		return page, nil, fmt.Errorf("can't extract page number from url %s in pagination link", pageUrl.String())
	}

	page, err = strconv.Atoi(pieces[len(pieces)-2])

	nextLink := link.Parent().Next().Find("a")
	if nextLink.Nodes != nil {
		nextHref, exists := nextLink.Attr("href")
		if !exists {
			err = fmt.Errorf("missing href from next page link")
			return
		}

		nextPageUrl, err = doc.Url.Parse(nextHref)
		if err != nil {
			err = fmt.Errorf("couldn't parse href for pagination next link")
			return
		}
	}

	return
}
