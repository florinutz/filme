package coll33tx

import (
	"net/url"

	"strconv"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

type ListItem struct {
	Name    string
	Href    string
	Size    string
	Seeds   int
	Leeches int
}

// TorrentFoundCallback is the type for a callback func to be called when the we came across a list item
type ListItemFoundCallback func(item ListItem)

// listCollector is a wrapper around the colly collector + listing page data
type listCollector struct {
	*colly.Collector
	onItemFound ListItemFoundCallback
}

func NewListCollector(onItemFound ListItemFoundCallback, options ...func(collector *colly.Collector)) *listCollector {
	listCollector := listCollector{
		Collector:   getCollyCollector(options...),
		onItemFound: onItemFound,
	}

	listCollector.Collector.OnHTML("td.name a:nth-of-type(2)", listCollector.itemHandler)

	return &listCollector
}

// onListItemHandler parses the  list and launches a request for each item page
func (col *listCollector) itemHandler(e *colly.HTMLElement) {
	linkHref := e.Request.AbsoluteURL(e.Attr("href"))

	log.WithFields(log.Fields{"title": e.Text, "href": linkHref}).Debug("list item")

	_, err := url.Parse(linkHref)
	if err != nil {
		log.WithError(err).WithField("href", linkHref).Warn("Invalid href in list, skipping it.")
		return
	}

	// trigger the onFound event with the data as input
	listItem := &ListItem{}

	errs := listItem.fromTitleLink(e)
	if errs != nil {
		log.WithField("errs", errs).Errorf("item '%s' is invalid", e.Text)
	}

	// trigger the event even on incomplete ListItem due to strtoint conversion errors
	col.onItemFound(*listItem)
}

func (i *ListItem) fromTitleLink(e *colly.HTMLElement) (errs []error) {
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
