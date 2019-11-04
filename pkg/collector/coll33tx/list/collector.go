package list

import (
	"fmt"

	"github.com/florinutz/filme/pkg/collector/coll33tx"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// PageCrawledCallback represents callback code that has access to all page data
type PageCrawledCallback func(lines []*Line, pagination *Pagination, r *colly.Response, log *log.Entry)

// ListCollector is a wrapper around the colly collector + listing page data
type Collector struct {
	*colly.Collector
	wantedItems   int
	pagesNeeded   int
	OnPageCrawled PageCrawledCallback
	Log           *log.Entry
}

// NewCollector instantiates a list page collector
func NewCollector(
	onPageCrawled PageCrawledCallback,
	wantedItems int,
	log *log.Entry,
	options ...func(collector *colly.Collector),
) *Collector {
	if wantedItems == 0 {
		panic("pretty sure you want more than 0 items")
	}

	c := colly.NewCollector(options...)

	coll33tx.DomainConfig(c, log)

	col := Collector{
		Collector:     c,
		wantedItems:   wantedItems,
		pagesNeeded:   0,
		OnPageCrawled: onPageCrawled,
		Log:           log,
	}

	col.pagesNeeded = (col.wantedItems-1)/ItemsPerPage + 1

	col.Collector.OnScraped(func(resp *colly.Response) {
		doc, err := NewDocument(resp)
		if err != nil {
			col.Log.WithError(err).Error("couldn't parse page document")
			return
		}

		lines := doc.GetLines()

		pagination := doc.GetPagination()
		col.OnPageCrawled(lines, pagination, resp, col.Log)

		if pagination != nil && pagination.Current == 1 && col.pagesNeeded > 1 {
			// This is the first page out of many, so let's launch parallel Visits to as many of them as we need to
			for pageNo := 2; pageNo <= col.pagesNeeded; pageNo++ {
				pUrl := fmt.Sprintf(pagination.LinksTpl, pageNo)

				if err := col.Visit(pUrl); err != nil {
					errMsg := fmt.Sprintf("couldn't initialize the visiting of page %d", pageNo)
					col.Log.WithError(err).Error(errMsg)
				}
			}
		}
	})

	return &col
}
