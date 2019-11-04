package filme

import (
	"encoding/json"
	"fmt"

	"github.com/florinutz/filme/pkg/collector/imdb/detail"
)

func (f *Filme) VisitImdbDetailPage(
	pageUrl string,
	jsonOutput bool,
) error {
	log := f.Log.WithField("url", pageUrl)

	func() {
		var onImdbItemFound = func(item *detail.ImdbFilm) {
			f.Log.WithField("item", item).Debug("imdb item found")

			if jsonOutput {
				j, err := json.Marshal(item)
				if err != nil {
					log.WithError(err).Fatal("error encoding to json")
				}
				fmt.Fprintln(f.Out, string(j))

				return
			}

			fmt.Fprintf(f.Out, `title: %s
Year: %d`,
				item.Title,
				item.Year,
			)
		}

		col := detail.NewCollector(onImdbItemFound, log)
		if err := col.Visit(pageUrl); err != nil {
			log.WithError(err).Fatal("visit error")
		}
		col.Wait()
	}()

	return nil
}
