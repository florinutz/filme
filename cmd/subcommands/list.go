package subcommands

import (
	"text/template"

	"os"

	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/phlo/filme/collector/coll33tx"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var (
	filmTemplate = template.New("film").Funcs(template.FuncMap{
		"h":       color.New(color.FgBlue).Sprint,
		"heading": color.New(color.FgBlue, color.Bold).Sprint,
		"highlightPrefix": func(prefixLength int, input string) string {
			return color.New(color.FgBlue).Sprint(input[:prefixLength]) + input[prefixLength:]
		},
		"highlightSubstr": func(what, input string) string {
			return strings.Replace(input, what, color.New(color.FgBlue).Sprint(what), -1)
		},
	})

	parsedTemplate = template.Must(filmTemplate.Parse(`

  {{ .Title | heading }}{{if .FilmTitle}}({{.FilmTitle|h}}{{if .FilmLink}} - {{.FilmLink|h}}{{end}}){{end}} - {{ .FoundOn | h }}
   {{if .IMDB}}{{ .IMDB.String|highlightSubstr "imdb"}}{{end}}
   {{ .Magnet | highlightPrefix 6 }}
   category {{ .Category | h }}, type {{ .Type | h }}, language {{ .Language | h }}, size {{ .TotalSize | h }}, downloads {{ .Downloads | h }}, seeders {{ .Seeds | h }}, leechers {{ .Leeches | h }}{{if .Image}}, image: {{.Image|h}}{{end}}
   {{if .FilmDescription}}description: {{.FilmDescription|h}}{{end}}`))

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
