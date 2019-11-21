package filme

import (
	"fmt"

	"github.com/florinutz/filme/pkg/collector/coll33tx/list"
	"github.com/florinutz/filme/pkg/config/value"
	"github.com/florinutz/filme/pkg/filme/l33tx_movies"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

func (f *Filme) Search(what string, requiredItems int, goIntoDetails bool, category value.LeetxListSearchCategory,
	movieEncoding value.LeetxListEncoding, sort value.LeetxListSortValue, debugLevel value.DebugLevelValue) error {
	log := f.Log.WithFields(map[string]interface{}{
		"search":   what,
		"category": category,
		"encoding": movieEncoding,
		"sort":     sort,
	})

	col := list.NewCollector(f.onSearchPageCrawled, requiredItems, *log)

	startUrl, err := l33tx_movies.GetListUrl(what, sort, &category, &movieEncoding)
	if err != nil {
		return fmt.Errorf("can't get start url: %w", err)
	}

	log.WithField("searchStartUrl", startUrl).Debug("starting search")

	if err = col.Visit(startUrl.String()); err != nil {
		log.WithError(err).Warn("initial search visit error")
		return err
	}

	col.Wait()

	return nil
}

func (f *Filme) onSearchPageCrawled(lines []*list.Line, pagination *list.Pagination, r *colly.Response, log logrus.Entry) {
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
