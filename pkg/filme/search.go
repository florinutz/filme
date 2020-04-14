package filme

import (
	"fmt"

	"github.com/florinutz/filme/pkg/filme/l33tx/list"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/filter"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/input"
	log "github.com/sirupsen/logrus"
)

type Searcher interface {
	Search(goIntoDetails bool, inputs input.ListingInput, filters filter.Filter,
		delay int, randomDelay int, parallelism int, userAgent string, log log.Entry) error
}

func (f *Filme) Search(goIntoDetails bool, inputs input.ListingInput, filters filter.Filter,
	delay int, randomDelay int, parallelism int, userAgent string, log log.Entry) error {
	ls := list.NewList(inputs, filters, log)

	col := list.NewCollector(ls, delay, randomDelay, parallelism, userAgent)

	startUrl, err := inputs.GetStartUrl()
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
