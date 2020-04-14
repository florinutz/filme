package list

import (
	"fmt"

	"github.com/florinutz/filme/pkg/collector/coll33tx"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/line"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/sirupsen/logrus"
)

// PageCrawledCallback represents callback code that has access to all page data
type PageCrawledCallback func(
	lines []*line.Line,
	pagination *Pagination,
	r *colly.Response,
)

// ListCollector is a wrapper around the colly collector + listing page data
type Collector struct {
	*colly.Collector
	OnPageCrawled PageCrawledCallback
	Log           logrus.Entry
}

// NewCollector instantiates a list page collector
func NewCollector(
	ls *Container,
	delay, randomDelay, parallelism int,
	userAgent string,
	options ...func(collector *colly.Collector),
) *Collector {
	c := colly.NewCollector(options...)

	coll33tx.DomainConfig(c, delay, randomDelay, parallelism, userAgent, ls.Log)

	col := Collector{
		Collector:     c,
		OnPageCrawled: ls.AddPage,
		Log:           ls.Log,
	}

	col.Collector.OnScraped(onScraped(col, ls))

	return &col
}

func onScraped(col Collector, ls *Container) func(resp *colly.Response) {
	return func(resp *colly.Response) {
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

		currentPage := 1

		pagination := doc.GetPagination()

		if pagination != nil {
			currentPage = pagination.Current

			if currentPage == 1 {
				ls.paging.limitLow = 1
				ls.paging.limitHigh = pagination.PagesCount
				if ls.paging.itemsPerPage, err = doc.CountItems(); err != nil {
					log.WithError(err).Fatal("could not count the value for items per page")
					return
				}

				ls.paging.pagesToCrawl = ls.paging.getNextPages(ls.Filters.MaxItems)

				q, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 1000})

				for _, pageNo := range ls.paging.pagesToCrawl {
					if pageNo == 1 { // already crawled
						continue
					}

					pageUrl := fmt.Sprintf(pagination.LinksTpl, pageNo)

					if err := q.AddURL(pageUrl); err != nil {
						log.WithError(err).WithField("url", pageUrl).Error("can't add url to queue")
					}

					if err := col.Visit(pageUrl); err != nil {
						errMsg := fmt.Sprintf("couldn't initialize the visiting of page %d", pageNo)
						log.WithError(err).Error(errMsg)
					}
				}

				q.Run(col.Collector)
			}

			// skip if current page is outside the filtered range (should be only page 1)
			if !ls.paging.pageIsValid(currentPage, ls.Filters.MaxItems) {
				log.WithFields(map[string]interface{}{
					"page":  currentPage,
					"range": ls.paging.pagesToCrawl,
				}).Debug("page out of range")
				return
			}
		}

		// extract lines
		lines, err := doc.GetLines()
		if err != nil {
			log.WithError(err).Warn()
			return
		}

		// perform callback on the bunch of lines extracted the valid page
		col.OnPageCrawled(lines, pagination, resp)
	}
}
