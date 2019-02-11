package filme

import (
	"net/url"

	"github.com/florinutz/filme/pkg/config/value"

	"github.com/florinutz/filme/pkg/collector/coll33tx"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

func (f *Filme) OnListPageCrawled(
	items []*coll33tx.Item,
	currentPage int,
	nextPage *url.URL,
	response *colly.Response,
	log *logrus.Entry) {
	f.Log.Printf("currentPage %d\nnextPage %s\nItems:\n%q\n", currentPage, nextPage, items)
}

func (f *Filme) VisitList(
	listUrl string,
	goIntoDetails bool,
	debugLevel value.DebugLevelValue,
) error {
	log := f.Log.WithField("start_url", listUrl)

	list := coll33tx.NewListCollector(f.OnListPageCrawled, log)

	err := list.Visit(listUrl)
	if err != nil {
		log.WithError(err).Warn("visit error")
		return err
	}

	list.Wait()

	return nil
}
