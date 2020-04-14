package detail

import (
	"net/url"

	"github.com/florinutz/filme/pkg/collector/coll33tx"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// OnTorrentFound is the type the callback func that's be called when a Torrent was onItemFound
type OnTorrentFound func(torrent *Torrent)

// DetailCollector is a wrapper around the colly collector + page data
type Collector struct {
	*colly.Collector
	onTorrentFound OnTorrentFound
	Torrent        Torrent // this will be filled in the events
	Log            log.Entry
}

func NewCollector(onTorrentFound OnTorrentFound, delay, randomDelay, parallelism int, userAgent string,
	log log.Entry, options ...func(*colly.Collector)) *Collector {
	c := colly.NewCollector(options...)

	coll33tx.DomainConfig(c, delay, randomDelay, parallelism, userAgent, log)

	col := Collector{
		Collector:      c,
		onTorrentFound: onTorrentFound,
		Torrent:        Torrent{},
		Log:            log,
	}

	col.Collector.OnScraped(col.OnScraped)

	return &col
}

// OnScraped assembles and collects the Torrent struct at the end
func (col *Collector) OnScraped(r *colly.Response) {
	doc, err := NewDocument(r, col.Log)
	if err != nil {
		col.Log.WithError(err).Error("couldn't parse page document")
	}
	col.onTorrentFound(doc.GetData())
}

// Torrent represents the data found on a detail page
type Torrent struct {
	ID              int
	Title           string
	FilmCleanTitle  string
	FilmLink        *url.URL
	FilmCategories  []string
	FilmDescription string
	IMDB            *url.URL
	Magnet          string
	FoundOn         *url.URL
	Image           *url.URL
	Description     string
	Genres          []string
	Category        string
	Type            string
	Language        string
	TotalSize       string
	UploadedBy      string
	Downloads       int
	LastChecked     string
	DateUploaded    string
	Seeders         int
	Leechers        int
	Year            int
	Quality         string
	TitleInfo       coll33tx.TitleInfo
}
