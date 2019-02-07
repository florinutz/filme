package parser

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/florinutz/filme/pkg/collector/coll33tx"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func getDocFromDesponse(r *colly.Response) (doc *goquery.Document, err error) {
	return goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
}

const (
	SelectorListPagination       = ".pagination"
	SelectorListPaginationActive = "li.active"
)

type listDocument struct {
	*goquery.Document
	response *colly.Response
}

func NewListDocument(r *colly.Response) (*listDocument, error) {
	doc, err := getDocFromDesponse(r)
	if err != nil {
		return nil, err
	}
	return &listDocument{Document: doc, response: r}, nil
}

func (doc *listDocument) GetItems() (items []*coll33tx.Item, err error) {
	selector := "td.name a:nth-of-type(2)"
	selection := doc.Find(selector)
	if selection.Nodes == nil {
		err = fmt.Errorf("couldn't select page items using selector '%s'", selector)
		return
	}
	for _, node := range selection.Nodes {
		item := coll33tx.Item{
			Name: node.Data,
			// todo href?
		}
		items = append(items, &item)
	}
	return
}

func (doc *listDocument) getPagination() (*goquery.Selection, error) {
	selection := doc.Find(SelectorListPagination)
	if selection.Nodes == nil {
		return nil, fmt.Errorf("pagination not found at selector '%s'", SelectorListPagination)
	}
	return selection, nil
}

func (doc *listDocument) GetCurrentPage() (page int, err error) {
	pagination, err := doc.getPagination()
	if err != nil {
		return
	}

	active := pagination.Find(SelectorListPaginationActive)
	if active.Nodes == nil {
		err = fmt.Errorf("pagination exists, but active page was not found in it (%s %s)",
			SelectorListPagination, SelectorListPaginationActive)
		return
	}

	return strconv.Atoi(active.Text())
}
