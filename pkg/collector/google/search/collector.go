package search

import (
	"github.com/florinutz/filme/pkg/collector/google"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// OnFound is the callback func type for when items were found
type OnFound func(items map[int]BaseItem, err error, onlyFilmRelatedItems bool, r *colly.Response, log *log.Entry)

// Collector is a wrapper around the colly collector + page data
type Collector struct {
	*colly.Collector
	onFound              OnFound
	Items                map[int]ItemDefault
	onlyFilmRelatedItems bool
	Log                  *log.Entry
}

func NewCollector(onFound OnFound, onlyFilmRelatedItems bool, log *log.Entry,
	options ...func(*colly.Collector)) *Collector {
	c := colly.NewCollector(options...)
	google.DomainConfig(c, log)

	col := Collector{
		Collector:            c,
		onFound:              onFound,
		onlyFilmRelatedItems: onlyFilmRelatedItems,
		Log:                  log,
	}

	col.Collector.OnScraped(col.onScraped)

	return &col
}

// onScraped assembles and collects the ImdbFilm struct at the end
func (col *Collector) onScraped(r *colly.Response) {
	doc, err := NewDocument(r, col.Log)
	if err != nil {
		col.Log.WithError(err).Error("couldn't parse detail page")
	}
	items, err := doc.GetItems()
	col.onFound(items, err, col.onlyFilmRelatedItems, r, col.Log)
}
