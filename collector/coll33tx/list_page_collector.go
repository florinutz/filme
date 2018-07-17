package coll33tx

import (
	"net/url"

	"strconv"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type Item struct {
	Name    string
	Href    string
	Size    string
	Seeds   int
	Leeches int
}

type List []*Item

// TorrentFoundCallbsack is the type for a callback func to be called when we came across a list item
type ListItemFoundCallback func(item Item)

//CrawlFinishedCallback is the type for the callback func to be called when all the list list were parsed into the list
type CrawlFinishedCallback func(list List, response *colly.Response)

// ListCollector is a wrapper around the colly collector + listing page data
type ListCollector struct {
	*colly.Collector
	list        List
	OnItemFound ListItemFoundCallback
	OnCrawled   CrawlFinishedCallback
}

func NewListCollector(
	onItemFound ListItemFoundCallback,
	onListCrawled CrawlFinishedCallback,
	options ...func(collector *colly.Collector),
) *ListCollector {
	listCollector := ListCollector{
		Collector:   getCollyCollector(options...),
		OnItemFound: onItemFound,
		OnCrawled:   onListCrawled,
	}
	listCollector.Collector.OnHTML("td.name a:nth-of-type(2)", listCollector.itemHandler)
	listCollector.Collector.OnScraped(func(r *colly.Response) {
		listCollector.OnCrawled(listCollector.list, r)
	})

	return &listCollector
}

// onListItemHandler parses the  list and launches a request for each item page
func (col *ListCollector) itemHandler(e *colly.HTMLElement) {
	linkHref := e.Request.AbsoluteURL(e.Attr("href"))

	log.WithFields(log.Fields{"title": e.Text, "href": linkHref}).Debug("list item")

	_, err := url.Parse(linkHref)
	if err != nil {
		log.WithError(err).WithField("href", linkHref).Warn("Invalid href in list, skipping it.")
		return
	}

	// trigger the onFound event with the data as input
	listItem := &Item{}

	errs := listItem.fromTitleLink(e)
	if errs != nil {
		log.WithField("errs", errs).Errorf("item '%s' is invalid", e.Text)
	}

	col.list = append(col.list, listItem)

	// trigger the event even on incomplete Item due to strtoint conversion errors
	col.OnItemFound(*listItem)
}

func (i *Item) fromTitleLink(e *colly.HTMLElement) (errs []error) {
	i.Name = e.Text
	i.Href = e.Request.AbsoluteURL(e.Attr("href"))

	tr := e.DOM.Parent().Parent()

	if seeds, err := strconv.Atoi(tr.Find(".seeds").Text()); err != nil {
		errs = append(errs, err)
	} else {
		i.Seeds = seeds
	}

	if leeches, err := strconv.Atoi(tr.Find(".leeches").Text()); err != nil {
		errs = append(errs, err)
	} else {
		i.Leeches = leeches
	}

	i.Size = tr.Find(".size").Text()

	return
}
