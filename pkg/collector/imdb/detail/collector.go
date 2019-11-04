package detail

import (
	"github.com/florinutz/filme/pkg/collector/imdb"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// OnFound is the type the callback func that's be called when a ImdbFilm was onItemFound
type OnFound func(*ImdbFilm)

// Collector is a wrapper around the colly collector + page data
type Collector struct {
	*colly.Collector
	onFound OnFound
	Movie   ImdbFilm // this will be filled in the events
	Log     *log.Entry
}

func NewCollector(onFound OnFound, log *log.Entry, options ...func(*colly.Collector)) *Collector {
	c := colly.NewCollector(options...)
	imdb.DomainConfig(c, log)

	col := Collector{
		Collector: c,
		onFound:   onFound,
		Movie:     ImdbFilm{},
		Log:       log,
	}

	col.Collector.OnScraped(col.OnScraped)

	return &col
}

// OnScraped assembles and collects the ImdbFilm struct at the end
func (col *Collector) OnScraped(r *colly.Response) {
	doc, err := NewDocument(r, col.Log)
	if err != nil {
		col.Log.WithError(err).Error("couldn't parse detail page")
	}
	col.onFound(doc.GetData())
}
