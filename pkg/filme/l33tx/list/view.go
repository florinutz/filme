package list

import (
	"html/template"
	"reflect"

	"github.com/fatih/color"
	"github.com/florinutz/filme/pkg/filme/l33tx/list/line"
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
	tpls = template.Must(template.New("").Funcs(funcMap).Parse(`{{ define "line" }}
{{.Item.Name | highlight}} ({{.Item.Size | bold}}, {{.Item.Seeders}} seeders, {{.Item.Leechers}} leechers): {{.Item.Href | bold}}
{{- range .Errs}}
    {{- .}}
{{end -}}
{{ end }}
`))
}

func DisplayLines(lines []line.Line) {
}
