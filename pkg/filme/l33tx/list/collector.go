package list

import (
	"fmt"
	"io"

	"github.com/florinutz/filme/pkg/collector/coll33tx"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/line"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

// PageCrawledCallback represents callback code that has access to all page data
type PageCrawledCallback func(
	w io.Writer,
	lines []*line.Line,
	pagination *Pagination,
	r *colly.Response,
	logger logrus.Entry,
) (itemsWritten uint)

// ListCollector is a wrapper around the colly collector + listing page data
type Collector struct {
	*colly.Collector
	wantedItems   uint
	pagesNeeded   uint
	OnPageCrawled PageCrawledCallback
	Log           logrus.Entry
}

// NewCollector instantiates a list page collector
func NewCollector(
	ls Container,
	options ...func(collector *colly.Collector),
) *Collector {
	c := colly.NewCollector(options...)

	coll33tx.DomainConfig(c, ls.Log)

	col := Collector{
		Collector:     c,
		wantedItems:   ls.Filters.MaxItems,
		pagesNeeded:   0,
		OnPageCrawled: ls.WritePage,
		Log:           ls.Log,
	}

	col.pagesNeeded = (col.wantedItems-1)/ItemsPerPage + 1

	col.Collector.OnScraped(func(resp *colly.Response) {
		log := col.Log.WithFields(map[string]interface{}{
			"url":    resp.Request.URL,
			"status": resp.StatusCode,
		})

		doc, err := NewDocument(resp)
		if err != nil {
			log.WithError(err).Error("problem creating collector")
			return
		}

		log = log.WithField("title", doc.Find("title").Text())

		lines, err := doc.GetLines()
		if err != nil {
			log.WithError(err).Warn()
			return
		}

		pagination := doc.GetPagination()
		col.OnPageCrawled(ls.Out, lines, pagination, resp, *log)

		if pagination != nil && pagination.Current == 1 && col.pagesNeeded > 1 {
			// This is the first page out of many, so let's launch parallel Visits to as many of them as we need to
			for pageNo := uint(2); pageNo <= col.pagesNeeded; pageNo++ {
				pUrl := fmt.Sprintf(pagination.LinksTpl, pageNo)

				if err := col.Visit(pUrl); err != nil {
					errMsg := fmt.Sprintf("couldn't initialize the visiting of page %d", pageNo)
					log.WithError(err).Error(errMsg)
				}
			}
		}
	})

	return &col
}
