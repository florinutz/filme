package filme

import (
	"fmt"

	"github.com/florinutz/filme/pkg/filme/l33tx/list"
)

type Searcher interface {
	Search(ls *list.Container, goIntoDetails bool, delay int, randomDelay int, parallelism int, userAgent string) error
}

func (f *Filme) Search(ls *list.Container, goIntoDetails bool, delay int, randomDelay int, parallelism int,
	userAgent string) error {
	col := list.NewCollector(ls, delay, randomDelay, parallelism, userAgent)

	startUrl, err := ls.Inputs.GetStartUrl()
	if err != nil {
		f.Log.WithError(err).Errorln()
		return fmt.Errorf("could not assemble the url: %w\n\n", err)
	}
	if startUrl == nil {
		f.Log.Fatal("empty url retrieved for search, please investigate")
	}

	f.Log.WithField("start_url", startUrl).Info("starting search")

	if err = col.Visit(startUrl.String()); err != nil {
		f.Log.WithError(err).Warn("initial search visit error")
		return err
	}

	col.Wait()

	ls.Display(f.Out)

	return nil
}
