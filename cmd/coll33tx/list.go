package coll33tx

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/florinutz/filme/util/value"

	"github.com/gocolly/colly"

	"github.com/florinutz/filme/collector/coll33tx"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"text/template"

	"github.com/fatih/color"
)

func init() {
	ListCmd.Flags().BoolVarP(&listCmdConfig.withDetails, "crawl-details", "d", false,
		"follows every link in the list and fetches more data")

	defaultDebugLevel := log.DebugLevel
	_ = listCmdConfig.debugLevel.Set(defaultDebugLevel.String())
	ListCmd.Flags().Var(&listCmdConfig.debugLevel, "debug-level", fmt.Sprintf("possible debug levels: %s",
		strings.Join(value.GetAllLevels(), ", ")))

	log.SetFormatter(&log.JSONFormatter{})

	L33txRootCmd.AddCommand(ListCmd)
}

type lCmdConfigType struct {
	withDetails bool
	debugLevel  value.DebugLevelValue
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

			list := coll33tx.NewListCollector(OnListItemFound, OnListCrawled, log)

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

	tplTorrentFull = template.Must(filmTemplate.Parse(`

  {{ .Title | heading }}{{if .FilmCleanTitle}}({{.FilmCleanTitle|h}}{{if .FilmLink}} - {{.FilmLink|h}}{{end}}){{end}} - {{ .FoundOn | h }}
   {{if .IMDB}}{{ .IMDB.String|highlightSubstr "imdb"}}{{end}}
   {{ .Magnet | highlightPrefix 6 }}
   category {{ .Category | h }}, type {{ .Type | h }}, language {{ .Language | h }}, size {{ .TotalSize | h }}, downloads {{ .Downloads | h }}, seeders {{ .Seeders | h }}, leechers {{ .Leechers | h }}{{if .Image}}, image: {{.Image|h}}{{end}}
   {{if .FilmDescription}}description: {{.FilmDescription|h}}{{end}}`))
)

// OnListItemFound is the callback executed when a new list item was found
func OnListItemFound(item coll33tx.Item, col *coll33tx.ListCollector, r *colly.Response) {
	logWithItem := log.WithField("item", item).WithField("url", r.Request.URL.String())

	if !listCmdConfig.withDetails {
		logWithItem.WithField("skip", true).Debug("list item found")
		return
	}

	// go deeper
	details := coll33tx.NewDetailsPageCollector(func(torrent coll33tx.Torrent) {
		if err := tplTorrentFull.Execute(os.Stdout, torrent); err != nil {
			logWithItem.WithError(err).Fatal("error while executing template")
		}
	}, col.Log)

	if err := details.Visit(item.Href); err != nil {
		logWithItem.WithError(err).Warn("visit error on list item page")
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
