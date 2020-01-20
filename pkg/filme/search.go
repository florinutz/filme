package filme

import (
	"fmt"

	"github.com/florinutz/filme/pkg/filme/l33tx/list"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/input"
)

func (f *Filme) Search(goIntoDetails bool, inputs input.ListingInput, filters filter.Filter) error {
	logFields := localLogFields(inputs, goIntoDetails, filters)
	log := *f.Log.WithFields(logFields)

	ls := list.NewList(inputs, filters, f.Out, log)

	col := list.NewCollector(ls)

	startUrl, err := ls.GetStartUrl()
	if err != nil {
		log.WithError(err).Errorln()
		return fmt.Errorf("could not assemble the url: %w\n\n", err)
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

	ls.Display(f.Out)

	return nil
}

func localLogFields(inputs input.ListingInput, goIntoDetails bool, filters filter.Filter) (result map[string]interface{}) {
	result = map[string]interface{}{
		"search_inputs_sort": inputs.Sort,
		"with_details":       goIntoDetails,
	}
	if inputs.Encoding != nil && *inputs.Encoding != 0 {
		result["search_inputs_encoding"] = *inputs.Encoding
	}
	if inputs.Category != nil && *inputs.Category != 0 {
		result["search_inputs_category"] = *inputs.Category
	}
	if inputs.Search != "" {
		result["search"] = inputs.Search
	}
	if inputs.URL != nil {
		result["search_inputs_url"] = inputs.URL.String()
	}

	for key, v := range filters.GetLogFields() {
		result[key] = v
	}

	return
}
