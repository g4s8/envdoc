package render

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	texttmpl "text/template"
)

//go:embed templ
var templatesFS embed.FS

var tplFuncs = map[string]any{
	// Standard string functions.
	// The functions here were added arbitrarily, and more can be added when needed.
	"repeat":     strings.Repeat,
	"split":      strings.Split,
	"join":       strings.Join,
	"contains":   strings.Contains,
	"toLower":    strings.ToLower,
	"toUpper":    strings.ToUpper,
	"toTitle":    strings.ToTitle,
	"replace":    strings.Replace,
	"hasPrefix":  strings.HasPrefix,
	"hasSuffix":  strings.HasSuffix,
	"trimSpace":  strings.TrimSpace,
	"trimPrefix": strings.TrimPrefix,
	"trimSuffix": strings.TrimSuffix,
	"trimLeft":   strings.TrimLeft,
	"trimRight":  strings.TrimRight,

	// Custom functions.
	"strAppend": func(arr []string, item string) []string {
		return append(arr, item)
	},
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
	"marshalIndent": func(v any) (string, error) {
		a, err := json.MarshalIndent(v, "", "  ")
		return string(a), err
	},
}

const (
	tmplDir     = "templ"
	tmplHelpers = "helpers.tmpl"
)

// newTmplText generates the template for a built-in format.
func newTmplText(name string) *texttmpl.Template {
	return texttmpl.Must(texttmpl.New(name).
		Funcs(tplFuncs).
		ParseFS(templatesFS,
			path.Join(tmplDir, name),
			path.Join(tmplDir, tmplHelpers)))
}

// newTmplCustom generates the template based on the provided custom template file.
func newTmplCustom(tmplFilePath string) (*texttmpl.Template, error) {
	tmpl, err := texttmpl.New(path.Base(tmplFilePath)).
		Funcs(tplFuncs).
		ParseFiles(tmplFilePath)
	if err != nil {
		return nil, fmt.Errorf("parsing template file: %w", err)
	}

	return tmpl, nil
}
