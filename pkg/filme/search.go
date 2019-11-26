package filme

import (
	"github.com/florinutz/filme/pkg/config/value"
	"github.com/florinutz/filme/pkg/filme/l33tx/list"
	filt "github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
	input2 "github.com/florinutz/filme/pkg/filme/l33tx/list/input"
)

func (f *Filme) Search(goIntoDetails bool, inputs input2.ListingInput, filters filt.Filter,
	debugLevel value.DebugLevelValue) error {

	log := *f.Log.WithFields(map[string]interface{}{
		"search_inputs": inputs,
		"with_details":  goIntoDetails,
		"filters":       filters,
	})

	ls := list.NewList(inputs, filters, f.Out, log)

	col := list.NewCollector(*ls)

	startUrl, err := ls.GetStartUrl()
	if startUrl == nil {
		log.Fatal("empty url retrieved for search, please investigate")
	}
	if err != nil {
		log.WithError(err).Fatal()
	}

	log.WithField("url", startUrl).Info("starting search")

	if err = col.Visit(startUrl.String()); err != nil {
		log.WithError(err).Warn("initial search visit error")
		return err
	}

	col.Wait()

	return nil
}
