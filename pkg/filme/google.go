package filme

import (
	"os"
	"reflect"
	"text/template"

	"github.com/fatih/color"
	"github.com/florinutz/filme/pkg/collector/google/search"
	"github.com/florinutz/filme/pkg/config/value"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

var tpls *template.Template

func init() {
	funcMap := template.FuncMap{
		"highlight": color.New(color.FgHiGreen).Add(color.Bold).SprintFunc(),
		"bold":      color.New(color.Bold).SprintFunc(),
		"avail": func(name string, data interface{}) bool {
			v := reflect.ValueOf(data)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			if v.Kind() != reflect.Struct {
				return false
			}
			return v.MethodByName(name).IsValid()
		},
	}
	tpls = template.Must(template.New("").Funcs(funcMap).Parse(`
{{ define "full_list" }}

{{- range $key, $value := . }}
{{ template "item" $value }}
{{ else }}
no results
{{ end -}}

{{ end }}
	
{{ define "item" -}}
{{.Title | highlight}}
{{.Url}}
{{.Description | bold}}
{{- range .Errs}}
    {{- .}}
{{end -}}
{{ end }}
	
{{ define "sources_list" }}{{- range $key, $items := . }}{{if $items -}}

{{$key | highlight}}:
{{ range $items }}{{ if ne $key "other" }}
{{- .Url }}
{{- if (avail "RatingStr" .) and (ne .RatingStr "")}} ({{.RatingStr}}){{end}}
{{- else -}}
{{- .Url }}
{{- end }}

{{end}}{{end}}{{end}}{{end}}
	`))
}

func (f *Filme) SearchGoogle(url string, onlyFilmRelatedItems bool, debugLevel value.DebugLevelValue) error {
	log := f.Log.WithField("start_url", url)
	col := search.NewCollector(f.onGooglePageCrawled, onlyFilmRelatedItems, log)
	err := col.Visit(url)
	if err != nil {
		log.WithError(err).Warn("visit error")
		return err
	}
	col.Wait()
	return nil
}

func (f *Filme) onGooglePageCrawled(items map[int]search.BaseItem, err error, onlyFilmRelatedItems bool,
	r *colly.Response, log *logrus.Entry) {
	if onlyFilmRelatedItems {
		sources := search.GetFilmSources(items)
		if err = tpls.ExecuteTemplate(os.Stdout, "sources_list", sources); err != nil {
			log.Printf("error rendering list:\n%s\n", err)
		}
		return
	}
	if err = tpls.ExecuteTemplate(os.Stdout, "full_list", items); err != nil {
		log.Printf("error rendering list:\n%s\n", err)
	}
}
