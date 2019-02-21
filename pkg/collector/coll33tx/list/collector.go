package list

import (
	"github.com/florinutz/filme/pkg/collector/coll33tx"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// PageCrawledCallback represents callback code that has access to all page data
type PageCrawledCallback func(lines []*Line, pagination *Pagination, r *colly.Response, log *log.Entry)

// ListCollector is a wrapper around the colly collector + listing page data
type Collector struct {
	*colly.Collector
	OnPageCrawled PageCrawledCallback
	Log           *log.Entry
}

// NewCollector instantiates a list page collector
func NewCollector(
	onPageCrawled PageCrawledCallback,
	log *log.Entry,
	options ...func(collector *colly.Collector),
) *Collector {
	c := colly.NewCollector(options...)

	coll33tx.DomainConfig(c, log)

	col := Collector{
		Collector:     c,
		OnPageCrawled: onPageCrawled,
		Log:           log,
	}

	col.Collector.OnScraped(func(r *colly.Response) {
		doc, err := NewDocument(r)
		if err != nil {
			col.Log.WithError(err).Error("couldn't parse page document")
			return
		}
		col.OnPageCrawled(doc.GetLines(), doc.GetPagination(), r, col.Log)
	})

	return &col
}
