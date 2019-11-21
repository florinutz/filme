package filme

import (
	"fmt"

	"github.com/florinutz/filme/pkg/collector/coll33tx/list"
	"github.com/florinutz/filme/pkg/config/value"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

func (f *Filme) Visit1337xListPage(
	listUrl string,
	goIntoDetails bool,
	debugLevel value.DebugLevelValue,
) error {
	log := f.Log.WithField("start_url", listUrl)

	col := list.NewCollector(f.On1337xListPageCrawled, 0, f.Out, f.Err, *log)

	err := col.Visit(listUrl)
	if err != nil {
		log.WithError(err).Warn("visit error")
		return err
	}

	col.Wait()

	return nil
}

func (f *Filme) On1337xListPageCrawled(lines []*list.Line, pagination *list.Pagination, r *colly.Response, log logrus.Entry) {
	if pagination != nil {
		fmt.Fprintf(f.Out, "current page: %d\n", pagination.Current)
		fmt.Fprintf(f.Out, "pages count: %d\n", pagination.PagesCount)
	}
	if len(lines) > 0 {
		fmt.Fprintln(f.Out, "")
	}
	for _, line := range lines {
		fmt.Fprintf(f.Out, "%s\n\t%s\n\tsize: %s, seeders: %d, leeches: %d\n\n",
			line.Item.Name,
			line.Item.Href,
			line.Item.Size,
			line.Item.Seeders,
			line.Item.Leeches)
		for _, err := range line.Errs {
			fmt.Fprintf(f.Err, "line error: %s", err)
		}
	}
}
