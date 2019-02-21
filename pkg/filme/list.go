package filme

import (
	"github.com/florinutz/filme/pkg/collector/coll33tx/list"
	"github.com/florinutz/filme/pkg/config/value"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

func (f *Filme) OnListPageCrawled(lines []*list.Line, pagination *list.Pagination, r *colly.Response, log *logrus.Entry) {
	if pagination != nil {
		f.Log.Printf("current page: %d\n", pagination.CurrentPage)
		f.Log.Printf("pages count: %d\n", pagination.PagesCount)
	}
	for _, line := range lines {
		f.Log.WithField("item", line.Item).Infoln()
		for _, err := range line.Errs {
			f.Log.Warnf("line error: %s", err)
		}
	}
}

func (f *Filme) VisitList(
	listUrl string,
	goIntoDetails bool,
	debugLevel value.DebugLevelValue,
) error {
	log := f.Log.WithField("start_url", listUrl)

	col := list.NewCollector(f.OnListPageCrawled, log)

	err := col.Visit(listUrl)
	if err != nil {
		log.WithError(err).Warn("visit error")
		return err
	}

	col.Wait()

	return nil
}
