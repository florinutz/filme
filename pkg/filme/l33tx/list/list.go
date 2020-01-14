package list

import (
	"fmt"
	"io"
	"net/url"
	"sort"

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
	Data         map[int][]*line.Line // stores pages
	maxItems     int
	itemsWritten int
	paging       *Paging
}

func NewList(inputs input.ListingInput, filters filter.Filter, out io.Writer, logger logrus.Entry) *Container {
	return &Container{
		Inputs:   inputs,
		Filters:  filters,
		Log:      logger,
		Out:      out,
		Data:     map[int][]*line.Line{},
		maxItems: int(filters.MaxItems),
		paging: &Paging{ // fill what's available here, while filling the rest when the first pagination is detected
			filterLow:  int(filters.Pages.Min),
			filterHigh: int(filters.Pages.Max),
		},
	}
}

func (l *Container) GetStartUrl() (*url.URL, error) {
	return l.Inputs.GetStartUrl()
}

func (l *Container) Display(w io.Writer) {
	var pages []int
	for pageNo, _ := range l.Data {
		pages = append(pages, pageNo)
	}

	sort.Ints(pages)

	for _, page := range pages {
		log := l.Log.WithField("page", page)

		log.Debug("displaying page")

		pageOutOfRange := !l.paging.pageIsValid(page, l.maxItems)
		if pageOutOfRange {
			log.WithField("range", l.paging.pagesToCrawl).Debug("page out of range, skipping it")
			break
		}

		for i, ln := range l.Data[page] {
			log = log.WithField("item", i)

			if errs := ln.Item.Validate(l.Filters); len(errs) > 0 {
				for _, err := range errs {
					log.WithError(err).Debug("item validation err")
				}
				continue
			}

			maxItemsReached := l.maxItems > 0 && l.itemsWritten >= l.maxItems
			if maxItemsReached {
				log.WithField("max", l.maxItems).Debug("max limit of items to display reached, stopping")
				break
			}

			fmt.Fprintf(w, "%d - %d: %s\n\t%s\n\tsize: %s, seeders: %d, leeches: %d\n\n",
				page,
				i+1,
				ln.Item.Name,
				ln.Item.Href,
				ln.Item.Size,
				ln.Item.Seeders,
				ln.Item.Leechers)

			for _, err := range ln.Errs {
				fmt.Fprintf(w, "line error: %s", err)
			}

			l.itemsWritten++
		}
	}
}

// AddPage will be called in loop for every new bunch of items retrieved
func (l *Container) AddPage(w io.Writer, lines []*line.Line, pagination *Pagination, r *colly.Response, logger logrus.Entry) {
	currentPage := 1
	if pagination != nil {
		currentPage = pagination.Current
	}

	l.Data[currentPage] = lines
}
