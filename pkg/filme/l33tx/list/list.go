package list

import (
	"fmt"
	"io"
	"net/url"

	"github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/input"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/line"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type Container struct {
	Inputs       input.ListingInput
	Filters      filter.Filter
	Out          io.Writer
	Log          logrus.Entry
	ItemsWritten int
}

func NewList(inputs input.ListingInput, filters filter.Filter, out io.Writer, logger logrus.Entry) *Container {
	return &Container{
		Inputs:  inputs,
		Filters: filters,
		Log:     logger,
		Out:     out,
	}
}

func (l *Container) GetStartUrl() (*url.URL, error) {
	return l.Inputs.GetStartUrl()
}

// WriteHeader will be called once at the beginning
func (l *Container) WriteHeader(w io.Writer) {
	//
}

// WritePage will be called in loop for every new bunch of items retrieved
func (l *Container) WritePage(w io.Writer, lines []*line.Line, pagination *Pagination, paging Paging,
	r *colly.Response, logger logrus.Entry) (itemsWritten int) {
	var currentPage int

	if pagination != nil {
		l.Log.WithFields(map[string]interface{}{
			"pagination_current": pagination.Current,
			"pagination_count":   pagination.PagesCount,
		}).Debug("pagination found")
		currentPage = pagination.Current
	} else {
		l.Log.Debugf("single page, no pagination")
		currentPage = 1
	}

	if len(lines) > 0 {
		fmt.Fprintln(w, "")
	}

	for i, ln := range lines {
		log := l.Log.WithField("item", ln.Item)
		if errs := ln.Item.Validate(l.Filters); len(errs) > 0 {
			for _, err := range errs {
				log.WithError(err).Debug("item validation err")
			}
			continue
		}

		currentItemOffset := i + 1 + l.ItemsWritten
		maxItemsReached := l.Filters.MaxItems > 0 && currentItemOffset > int(l.Filters.MaxItems)
		if maxItemsReached {
			log.WithField("max", l.Filters.MaxItems).Debug("max limit of items to display reached, stopping")
			break
		}

		pageOutOfRange := !paging.pageIsValid(currentPage, int(l.Filters.MaxItems))
		if pageOutOfRange {
			log.WithFields(map[string]interface{}{
				"page":  currentPage,
				"range": paging.pagesToCrawl,
			}).Debug("page out of range, skipping it")
			break
		}

		fmt.Fprintf(w, "%d: %s\n\t%s\n\tsize: %s, seeders: %d, leeches: %d\n\n",
			currentPage,
			ln.Item.Name,
			ln.Item.Href,
			ln.Item.Size,
			ln.Item.Seeders,
			ln.Item.Leechers)

		for _, err := range ln.Errs {
			fmt.Fprintf(w, "line error: %s", err)
		}

		itemsWritten++
	}

	return
}

// WriteFooter will be called once at the end
func (l *Container) WriteFooter(w io.Writer) {
	//
}
