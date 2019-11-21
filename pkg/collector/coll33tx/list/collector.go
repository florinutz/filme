package list

import (
	"fmt"
	"io"

	"github.com/florinutz/filme/pkg/collector/coll33tx"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

// PageCrawledCallback represents callback code that has access to all page data
type PageCrawledCallback func(lines []*Line, pagination *Pagination, wantedItems int, r *colly.Response, log logrus.Entry)

// ListCollector is a wrapper around the colly collector + listing page data
type Collector struct {
	*colly.Collector
	wantedItems   int
	pagesNeeded   int
	OnPageCrawled PageCrawledCallback
	Out           io.Writer
	Err           io.Writer
	Log           logrus.Entry
}

const LeetxItemsPerPage = 50

// NewCollector instantiates a list page collector
func NewCollector(
	onPageCrawled PageCrawledCallback,
	wantedItems int,
	out io.Writer,
	err io.Writer,
	log logrus.Entry,
	options ...func(collector *colly.Collector),
) *Collector {
	if wantedItems == 0 {
		panic("I'm pretty sure you want more than 0 items")
	}

	c := colly.NewCollector(options...)

	coll33tx.DomainConfig(c, log)

	col := Collector{
		Collector:     c,
		wantedItems:   wantedItems,
		pagesNeeded:   0,
		OnPageCrawled: onPageCrawled,
		Out:           out,
		Err:           err,
		Log:           log,
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
			fmt.Fprintf(col.Out, "%s\n", err)

			return
		}

		pagination := doc.GetPagination()
		col.OnPageCrawled(lines, pagination, wantedItems, resp, *log)

		if pagination != nil && pagination.Current == 1 && col.pagesNeeded > 1 {
			// This is the first page out of many, so let's launch parallel Visits to as many of them as we need to
			for pageNo := 2; pageNo <= col.pagesNeeded; pageNo++ {
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
