package filme

import (
	"fmt"

	"github.com/florinutz/filme/pkg/filme/l33tx/list"
	filt "github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
	input2 "github.com/florinutz/filme/pkg/filme/l33tx/list/input"
)

func (f *Filme) Search(goIntoDetails bool, inputs input2.ListingInput, filters filt.Filter) error {
	log := *f.Log.WithFields(map[string]interface{}{
		"search_inputs_encoding": inputs.Encoding,
		"search_inputs_category": inputs.Category,
		"search_inputs_sort":     inputs.Sort,
		"search_inputs_search":   inputs.Search,
		"search_inputs_url":      inputs.Url,
		"with_details":           goIntoDetails,
		"filter_max_items":       fmt.Sprint(filters.MaxItems),
		"filter_leechers":        fmt.Sprint(filters.Leechers),
		"filter_seeders":         fmt.Sprint(filters.Seeders),
		"filter_size":            fmt.Sprint(filters.Size),
	})

	ls := list.NewList(inputs, filters, f.Out, log)

	col := list.NewCollector(*ls)

	startUrl, err := ls.GetStartUrl()
	if err != nil {
		log.WithError(err).Fatal()
	}
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
