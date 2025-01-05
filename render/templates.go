package render

import (
	"embed"
	"encoding/json"
	"path"
	"strings"

	texttmpl "text/template"
)

//go:embed templ
var templatesFS embed.FS

var tplFuncs = map[string]any{
	"repeat": strings.Repeat,
	"split":  strings.Split,
	"strAppend": func(arr []string, item string) []string {
		return append(arr, item)
	},
	"join": strings.Join,
	"strSlice": func() []string {
		return make([]string, 0)
	},
	"list": func(args ...any) []any {
		return args
	},
	"sum": func(args ...int) int {
		var sum int
		for _, v := range args {
			sum += v
		}
		return sum
	},
	"marshal": func(v any) string {
		a, _ := json.Marshal(v)
		return string(a)
	},
}

const (
	tmplDir     = "templ"
	tmplHelpers = "helpers.tmpl"
)

func newTmplText(name string) *texttmpl.Template {
	return texttmpl.Must(texttmpl.New(name).
		Funcs(tplFuncs).
		ParseFS(templatesFS,
			path.Join(tmplDir, name),
			path.Join(tmplDir, tmplHelpers)))
}
