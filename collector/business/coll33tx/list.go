package coll33tx

import (
	"os"

	"text/template"

	"github.com/florinutz/filme/collector/coll33tx"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"

	"strings"

	"github.com/fatih/color"
)

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
)

// NewListCollector creates a new 1337x.to list collector tweaked for business
func NewListCollector(
	withDetails bool,
	log *logrus.Entry,
	options ...func(collector *colly.Collector),
) *coll33tx.ListCollector {
	return coll33tx.NewListCollector(
		OnListItemFound(withDetails, log),
		OnListCrawled(withDetails, log),
		options...,
	)
}

// OnListItemFound is the callback executed when a new list item was found
func OnListItemFound(crawlDetails bool, log *logrus.Entry) func(coll33tx.Item) {
	return func(item coll33tx.Item) {
		const listItemFoundStr = "list item found"

		itemLog := log.WithField("item", item)
		if !crawlDetails {
			itemLog.WithField("skip", true).Debug(listItemFoundStr)
			return
		}
		itemLog.Debug(listItemFoundStr)

		// go deeper
		details := coll33tx.NewDetailsPageCollector(func(torrent coll33tx.L33tTorrent) {
			parsedTemplate.Execute(os.Stdout, torrent)
		})

		err := details.Visit(item.Href)
		if err != nil {
			log.WithError(err).Warn("visit error")
		}

		details.Wait()
	}
}

// OnListCrawled is executed when all the list was finally parsed
func OnListCrawled(crawlDetails bool, log *logrus.Entry) func(coll33tx.List, *colly.Response) {
	return func(list coll33tx.List, r *colly.Response) {
		const listWasCrawled = "list was crawled"

		l := log.WithField("list", list)

		if crawlDetails {
			l.Debug(listWasCrawled)
		} else {
			l.Info(listWasCrawled)
		}
	}
}
