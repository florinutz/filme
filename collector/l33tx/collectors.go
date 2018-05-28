package l33tx

import (
	"net/url"

	"strings"

	"sync"

	"github.com/gocolly/colly"
	collyExtensions "github.com/gocolly/colly/extensions"
	log "github.com/sirupsen/logrus"
	"gitlab.com/phlo/filme/collector"
)

const Domain = "1337x.to"

var (
	ListCollector    *colly.Collector
	DetailsCollector *colly.Collector
	Torrents         collector.Torrents
	mutex            *sync.Mutex
)

func init() {
	ListCollector = NewListCollector()
	DetailsCollector = NewDetailsCollector(ListCollector)
	mutex = &sync.Mutex{}
}

func NewListCollector() *colly.Collector {
	c := colly.NewCollector()
	collector.Configure(c)
	c.AllowedDomains = []string{Domain}
	collyExtensions.RandomUserAgent(c)
	collyExtensions.Referrer(c)
	c.OnHTML("td.name a:nth-of-type(2)", listItemHandler)

	return c
}

func NewDetailsCollector(listCollector *colly.Collector) *colly.Collector {
	detailsCollector := listCollector.Clone()
	detailsCollector.OnHTML(".box-info-heading h1", titleLookup)
	detailsCollector.OnHTML("a[href]", magnetLookup)
	detailsCollector.OnResponse(responseHandler)
	detailsCollector.OnScraped(onScraped)

	return detailsCollector
}

func responseHandler(r *colly.Response) {
	mutex.Lock()
}

// onListItemHandler parses the  list and launches a request for each item page
func listItemHandler(e *colly.HTMLElement) {
	linkHref := e.Request.AbsoluteURL(e.Attr("href"))

	log.WithFields(log.Fields{"title": e.Text, "href": linkHref}).Debug("list item")

	_, err := url.Parse(linkHref)
	if err != nil {
		log.WithError(err).WithField("href", linkHref).Warn("Invalid href in list, skipping it.")
		return
	}

	// go deeper
	err = DetailsCollector.Request(e.Request.Method, linkHref, nil, e.Response.Ctx, nil)
	if err != nil {
		log.WithError(err).Warn("Visit error")
	}
}

// titleLookup gets the title from the details page
func titleLookup(e *colly.HTMLElement) {
	e.Response.Ctx.Put(collector.CtxKeyTitle, e.Text)
}

// magnetLookup gets the magnet link from the details page
func magnetLookup(e *colly.HTMLElement) {
	if !strings.HasPrefix(e.Attr("href"), "magnet") {
		return
	}

	e.Response.Ctx.Put(collector.CtxKeyMagnet, e.Attr("href"))
}

// onScraped assembles and collects the Torrent struct at the end
func onScraped(r *colly.Response) {
	defer mutex.Unlock()

	loggerWithURL := log.WithField("url", r.Request.URL)

	title := r.Ctx.Get(collector.CtxKeyTitle)
	if title == "" {
		loggerWithURL.Warn("The title was not harvested properly")
	}

	magnet := r.Ctx.Get(collector.CtxKeyMagnet)
	if magnet == "" {
		loggerWithURL.Warn("The magnet was not harvested properly")
	}

	Torrents = append(Torrents, &collector.Torrent{
		Title:   title,
		Magnet:  magnet,
		FoundOn: r.Request.URL,
	})
}
