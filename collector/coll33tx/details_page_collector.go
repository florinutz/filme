package coll33tx

import (
	"net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// L33tTorrent represents the data onItemFound on a torrent details page
type L33tTorrent struct {
	Magnet       string
	FoundOn      *url.URL
	Title        string
	Genres       []string
	Description  string
	Category     string
	Type         string
	Language     string
	TotalSize    string
	Downloads    int
	LastChecked  string
	DateUploaded string
	Seeders      int
	Leechers     int
}

// TorrentFoundCallback is the type the callback func that's be called when a torrent was onItemFound
type TorrentFoundCallback func(torrent L33tTorrent)

// detailsCollector is a wrapper around the colly collector + page data
type detailsCollector struct {
	*colly.Collector
	found   TorrentFoundCallback
	torrent L33tTorrent // this will be filled in the events
}

var (
	mutex *sync.Mutex
)

func init() {
	mutex = &sync.Mutex{}
}

func NewDetailsPageCollector(found TorrentFoundCallback, options ...func(*colly.Collector)) *detailsCollector {
	col := detailsCollector{
		Collector: getCollyCollector(options...),
		found:     found,
		torrent:   L33tTorrent{},
	}

	col.Collector.OnHTML(".box-info-heading h1", col.Title)
	col.Collector.OnHTML("a[href]", col.Magnet)
	col.Collector.OnResponse(col.OnResponse)
	col.Collector.OnScraped(col.OnScraped)

	return &col
}

// titleLookup gets the title from the details page
func (c *detailsCollector) Title(e *colly.HTMLElement) {
	c.torrent.Title = e.Text
}

// magnetLookup gets the magnet link from the details page
func (c *detailsCollector) Magnet(e *colly.HTMLElement) {
	if !strings.HasPrefix(e.Attr("href"), "magnet") {
		return
	}
	c.torrent.Magnet = e.Attr("href")
}

func (c *detailsCollector) OnResponse(r *colly.Response) {
	mutex.Lock()
}

// OnScraped assembles and collects the Torrent struct at the end
func (c *detailsCollector) OnScraped(r *colly.Response) {
	defer mutex.Unlock()

	loggerWithURL := log.WithField("url", r.Request.URL)

	if c.torrent.Title == "" {
		loggerWithURL.Warn("The title was not harvested properly")
	}

	if c.torrent.Magnet == "" {
		loggerWithURL.Warn("The magnet was not harvested properly")
	}

	if c.torrent.Category == "" {
		loggerWithURL.Warn("The category was not harvested properly")
	}

	c.found(c.torrent)
}
