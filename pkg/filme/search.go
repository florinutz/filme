package filme

import (
	"fmt"

	"github.com/florinutz/filme/pkg/collector/coll33tx/list"
	"github.com/florinutz/filme/pkg/config/value"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/encoding"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/search_category"
	"github.com/florinutz/filme/pkg/config/value/1337x/list/sort"
	filt "github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/url"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

func (f *Filme) Search(
	searchStr string,
	goIntoDetails bool,
	category search_category.SearchCategory,
	movieEncoding encoding.ListEncoding,
	sort sort.Value,
	filters filt.Filter,
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

	col := list.NewCollector(f.onSearchPageCrawled, filters, f.Out, f.Err, *log)

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
	clientSideFiltering filt.Filter,
	pagination *list.Pagination,
	r *colly.Response,
	log logrus.Entry,
) {
	currentPage := 1

	if pagination != nil {
		log.WithFields(map[string]interface{}{
			"pagination_current": pagination.Current,
			"pagination_count":   pagination.PagesCount,
		}).Debug("pagination found")
		currentPage = pagination.Current
	}

	if len(lines) > 0 {
		fmt.Fprintln(f.Out, "")
	}

	for i, line := range lines {
		log = *log.WithField("item", line.Item)
		if errs := line.Item.Validate(clientSideFiltering); len(errs) > 0 {
			for _, err := range errs {
				log.WithError(err).Debug("item validation err")
			}
			continue
		}
		currentItemOffset := i + 1 + (currentPage-1)*list.LeetxItemsPerPage
		if clientSideFiltering.MaxItems > 0 && currentItemOffset > int(clientSideFiltering.MaxItems) {
			log.WithField("max", clientSideFiltering.MaxItems).Debug("max limit of items to display reached, stopping")
			break
		}
		fmt.Fprintf(f.Out, "%s\n\t%s\n\tsize: %s, seeders: %d, leeches: %d\n\n",
			line.Item.Name,
			line.Item.Href,
			line.Item.Size,
			line.Item.Seeders,
			line.Item.Leechers)
		for _, err := range line.Errs {
			fmt.Fprintf(f.Err, "line error: %s", err)
		}
	}
}
