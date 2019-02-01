package coll33tx

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gocolly/colly"

	"github.com/florinutz/filme/collector/coll33tx"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"text/template"

	"github.com/fatih/color"
)

func init() {
	ListCmd.Flags().BoolVarP(&listCmdConfig.withDetails, "crawl-details", "d", false, "follows every link in the list and fetches more data")
	L33txRootCmd.AddCommand(ListCmd)
	log.SetFormatter(&log.JSONFormatter{})
}

type lCmdConfigType struct {
	withDetails bool
	url         string
}

var (
	listCmdConfig lCmdConfigType

	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "Parses 1337x listings",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if strings.Contains(args[0], "//1337x.to") {
					listCmdConfig.url = args[0]
				} else { // search
					params := url.Values{"search": {strings.Join(args, " ")}}
					listCmdConfig.url = fmt.Sprintf("https://1337x.to/srch?%s", params.Encode())
				}
			} else {
				listCmdConfig.url = "https://1337x.to/popular-movies"
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			log := log.WithField("url", listCmdConfig.url)

			list := coll33tx.NewListCollector(OnListItemFound, OnListCrawled)

			err := list.Visit(listCmdConfig.url)
			if err != nil {
				log.WithError(err).Warn("visit error")
			}
			list.Wait()
		},
	}

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
)

// OnListItemFound is the callback executed when a new list item was found
func OnListItemFound(item coll33tx.Item) {
	logWithItem := log.WithField("item", item)

	if !listCmdConfig.withDetails {
		logWithItem.WithField("skip", true).Debug("list item found")
		return
	}

	// go deeper
	details := coll33tx.NewDetailsPageCollector(func(torrent coll33tx.L33tTorrent) {
		if err := parsedTemplate.Execute(os.Stdout, torrent); err != nil {
			logWithItem.WithError(err).Fatal("error while executing template")
		}
	})

	if err := details.Visit(item.Href); err != nil {
		logWithItem.WithError(err).Fatal("visit error")
	}

	details.Wait()
}

// OnListCrawled is executed when all the list was finally parsed
func OnListCrawled(list coll33tx.List, r *colly.Response) {
	const listWasCrawled = "list was crawled"
	l := log.WithField("url", r.Request.URL.String())
	if listCmdConfig.withDetails {
		l.Debug(listWasCrawled)
	} else {
		l.Info(listWasCrawled)
	}

	for _, item := range list {
		fmt.Printf("%s: %s\n", item.Name, item.Href)
	}
}
