package filme

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eefret/gomdb"
	"github.com/florinutz/filme/pkg/collector/coll33tx/detail"
)

func (f *Filme) Visit1337xDetailPage(
	pageUrl string,
	justMagnet bool,
	jsonOutput bool,
) error {
	log := f.Log.WithField("url", pageUrl)

	func() {
		var OnTorrentFound detail.OnTorrentFound = func(torrent *detail.Torrent) {
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

id: %d
magnet: %s

seeders: %d
leechers: %d`,
				strings.Trim(torrent.Title, " "),
				torrent.ID,
				torrent.Magnet,
				torrent.Seeders,
				torrent.Leechers,
			)

			if omdbApiKey, ok := os.LookupEnv("OMDB_API_KEY"); ok {
				gomdbApi := gomdb.Init(omdbApiKey)
				query := &gomdb.QueryData{Title: torrent.FilmCleanTitle, SearchType: gomdb.MovieSearch}
				if torrent.Year > 0 {
					query.Year = strconv.Itoa(torrent.Year)
				}
				res, err := gomdbApi.Search(query)
				if err != nil {
					log.WithError(err).WithField("title", torrent.Title).Fatal("omdb lookup failed")
				}
				if len(res.Search) > 0 {
					fmt.Fprintf(f.Out, "\n\nImdb info:\n")
					for _, searchResult := range res.Search {
						imdbLink := fmt.Sprintf("https://www.imdb.com/title/%s/", searchResult.ImdbID)
						fmt.Fprintf(f.Out, "%s (%s) %s\n", imdbLink, searchResult.Year, searchResult.Title)
					}
				}
			} else {
				log.Warn("no omdb api key")
			}

		}

		col := detail.NewCollector(OnTorrentFound, *log)
		if err := col.Visit(pageUrl); err != nil {
			log.WithError(err).Fatal("visit error")
		}

		col.Wait()
	}()

	return nil
}
