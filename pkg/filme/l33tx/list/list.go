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

const ListItemsPerPage = 50

type Container struct {
	Inputs  input.ListingInput
	Filters filter.Filter
	Out     io.Writer
	Log     logrus.Entry
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
func (l *Container) WritePage(
	w io.Writer,
	lines []*line.Line,
	pagination *Pagination,
	r *colly.Response,
	logger logrus.Entry,
) (itemsWritten uint) {
	currentPage := 1

	if pagination != nil {
		l.Log.WithFields(map[string]interface{}{
			"pagination_current": pagination.Current,
			"pagination_count":   pagination.PagesCount,
		}).Debug("pagination found")
		currentPage = pagination.Current
	}

	if len(lines) > 0 {
		fmt.Fprintln(w, "")
	}

	for i, line := range lines {
		log := l.Log.WithField("item", line.Item)
		if errs := line.Item.Validate(l.Filters); len(errs) > 0 {
			for _, err := range errs {
				log.WithError(err).Debug("item validation err")
			}
			continue
		}

		currentItemOffset := i + 1 + (currentPage-1)*ListItemsPerPage
		if l.Filters.MaxItems > 0 && currentItemOffset > int(l.Filters.MaxItems) {
			log.WithField("max", l.Filters.MaxItems).Debug("max limit of items to display reached, stopping")
			break
		}

		fmt.Fprintf(w, "%s\n\t%s\n\tsize: %s, seeders: %d, leeches: %d\n\n",
			line.Item.Name,
			line.Item.Href,
			line.Item.Size,
			line.Item.Seeders,
			line.Item.Leechers)

		for _, err := range line.Errs {
			fmt.Fprintf(w, "line error: %s", err)
		}
	}

	return 0
}

// WriteFooter will be called once at the end
func (l *Container) WriteFooter(w io.Writer) {
	//
}
