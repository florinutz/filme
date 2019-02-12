package filme

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/eefret/gomdb"

	"github.com/florinutz/filme/pkg/collector/coll33tx"
)

func (f *Filme) VisitDetailPage(
	pageUrl string,
	justMagnet bool,
	jsonOutput bool,
) error {
	log := f.Log.WithField("url", pageUrl)

	func() {
		OnTorrentFound := func(torrent coll33tx.Torrent) {
			f.Log.WithField("torrent", torrent).Debug("torrent found on detail page")

			if justMagnet {
				fmt.Fprintln(f.Out, torrent.Magnet)

				return
			}

			if jsonOutput {
				j, err := json.Marshal(torrent)
				if err != nil {
					log.WithError(err).Fatal("error encoding to json")
				}
				fmt.Fprintln(f.Out, string(j))
				return
			}

			fmt.Fprintf(f.Out, `%s

magnet: %s

seeders: %d
leechers: %d`,
				strings.Trim(torrent.Title, " "),
				torrent.Magnet,
				torrent.Seeders,
				torrent.Leechers,
			)

			if omdbApiKey, ok := os.LookupEnv("OMDB_API_KEY"); ok {
				gomdbApi := gomdb.Init(omdbApiKey)
				query := &gomdb.QueryData{Title: torrent.FilmCleanTitle, SearchType: gomdb.MovieSearch}
				res, err := gomdbApi.Search(query)
				if err != nil {
					log.WithError(err).WithField("title", torrent.Title).Fatal("omdb lookup failed")
				}
				fmt.Fprintln(f.Out, res.Search)
			} else {
				log.Warn("no omdb api key")
			}

		}

		detail := coll33tx.NewDetailsPageCollector(OnTorrentFound, log)
		if err := detail.Visit(pageUrl); err != nil {
			log.WithError(err).Fatal("visit error")
		}

		detail.Wait()
	}()

	return nil
}
