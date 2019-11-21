package filme

import (
	"fmt"

	"github.com/florinutz/filme/pkg/collector/coll33tx/list"
	"github.com/florinutz/filme/pkg/config/value"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/encoding"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/search_category"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/sort"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/url"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

func (f *Filme) Search(
	searchStr string,
	requiredItems int,
	goIntoDetails bool,
	category search_category.SearchCategory,
	movieEncoding encoding.ListEncoding,
	sort sort.Value,
	debugLevel value.DebugLevelValue) error {

	startUrl, err := url.GetListUrl(searchStr, sort, &category, &movieEncoding)
	if err != nil {
		return fmt.Errorf("can't get start url: %w", err)
	}
	log := f.Log.WithFields(map[string]interface{}{
		"search":   searchStr,
		"category": category.String(),
		"encoding": movieEncoding.String(),
		"sort":     sort.String(),
	})

	col := list.NewCollector(f.onSearchPageCrawled, requiredItems, f.Out, f.Err, *log)

	if startUrl == nil {
		log.Fatal("empty url retrieved for search, please investigate")
	}

	log.WithField("url", startUrl).Info("starting search")

	if err = col.Visit(startUrl.String()); err != nil {
		log.WithError(err).Warn("initial search visit error")
		return err
	}

	col.Wait()

	return nil
}

func (f *Filme) onSearchPageCrawled(
	lines []*list.Line,
	pagination *list.Pagination,
	wantedItems int,
	r *colly.Response,
	log logrus.Entry,
) {
	var currentPage int

	if pagination != nil {
		log.WithFields(map[string]interface{}{
			"pagination_current": pagination.Current,
			"pagination_count":   pagination.PagesCount,
		}).Debug("pagination found")
		currentPage = pagination.Current
	} else {
		currentPage = 1
	}

	if len(lines) > 0 {
		fmt.Fprintln(f.Out, "")
	}

	for i, line := range lines {
		if i+1+currentPage*list.LeetxItemsPerPage > wantedItems {
			log.WithField("max", wantedItems).Debug("max limit of items to display reached, stopping")
			break
		}
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
