package coll33tx

import (
	"net/url"
	"strings"
	"sync"

	"bytes"

	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// L33tTorrent represents the data onItemFound on a torrent details page
type L33tTorrent struct {
	Title        string
	Magnet       string
	FoundOn      *url.URL
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
	Leeches      int
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

	col.Collector.OnResponse(col.OnResponse)
	col.Collector.OnScraped(col.OnScraped)

	return &col
}

// Magnet gets the magnet link from the details page
func (c *detailsCollector) Magnet(e *colly.HTMLElement) {
	if !strings.HasPrefix(e.Attr("href"), "magnet") {
		return
	}
	c.torrent.Magnet = e.Attr("href")
}

func (c *detailsCollector) OnResponse(r *colly.Response) {
	c.torrent.fromResponse(r)
}

// OnScraped assembles and collects the Torrent struct at the end
func (c *detailsCollector) OnScraped(r *colly.Response) {
	c.found(c.torrent)
}

func (torrent *L33tTorrent) fromResponse(r *colly.Response) (errs []error) {
	mutex.Lock()
	defer mutex.Unlock()

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(r.Body))
	if err != nil {
		return []error{err}
	}

	if title := doc.Find(".box-info-heading h1"); title.Nodes == nil {
		errs = append(errs, errors.New("missing title"))
	} else {
		torrent.Title = title.Text()
	}

	if links := doc.Find("a[href]").FilterFunction(func(_ int, s *goquery.Selection) bool {
		href, _ := s.Attr("href")
		return strings.HasPrefix(href, "magnet:?")
	}); links.Nodes == nil {
		errs = append(errs, errors.New("missing magnet"))
	} else {
		torrent.Magnet, _ = links.Attr("href")
	}

	return
}
