package subcommands

import (
	"text/template"

	"os"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/phlo/filme/collector/coll33tx"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var (
	highlight = color.New(color.FgBlue)

	filmTemplate = template.New("film").Funcs(template.FuncMap{
		"highlight": highlight.Sprint,
		"highlightPrefix": func(prefixLength int, input string) string {
			return highlight.Sprint(input[:prefixLength]) + input[prefixLength:]
		},
	})

	parsedTemplate = template.Must(filmTemplate.Parse(`
{{ .Title | highlight }}
  {{ .Magnet | highlightPrefix 6 }}
`))

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists l33t films",
		Long:  `Handles 1337x.to listings and torrent detail pages.`,

		Run: func(cmd *cobra.Command, args []string) {
			list := coll33tx.NewListCollector(onListItemFound)
			err := list.Visit("https://1337x.to/popular-movies")
			if err != nil {
				log.WithError(err).Warn("Visit error")
			}
			list.Wait()
		},
	}
)

func onListItemFound(item coll33tx.ListItem) {
	log.WithField("item", item).Debug("list item found")
	details := coll33tx.NewDetailsPageCollector(onTorrentFound)
	err := details.Visit(item.Href)
	if err != nil {
		log.WithError(err).Warn("Visit error")
	}
	details.Wait()
}

func onTorrentFound(torrent coll33tx.L33tTorrent) {
	parsedTemplate.Execute(os.Stdout, torrent)
}
