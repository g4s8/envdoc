package main

import (
	"embed"
	"fmt"
	"io"
	"strings"

	htmltmpl "html/template"
	texttmpl "text/template"
)

type renderItem struct {
	EnvName      string
	Doc          string
	EnvDefault   string
	EnvSeparator string

	Required bool
	Expand   bool
	NonEmpty bool
	FromFile bool
}

type renderContext struct {
	Items []renderItem
}

func newRenderContext(items []docItem) renderContext {
	res := renderContext{
		Items: make([]renderItem, len(items)),
	}
	for i, item := range items {
		res.Items[i] = renderItem{
			EnvName:      item.envName,
			Doc:          item.doc,
			EnvDefault:   item.envDefault,
			EnvSeparator: item.separator,
			Required:     item.flags&docItemFlagRequired != 0,
			Expand:       item.flags&docItemFlagExpand != 0,
			NonEmpty:     item.flags&docItemFlagNonEmpty != 0,
			FromFile:     item.flags&docItemFlagFromFile != 0,
		}
	}
	return res
}

//go:embed templ
var templates embed.FS

var renderFns = map[string]any{
	"join": strings.Join,
	"appendS": func(s []string, v string) []string {
		return append(s, v)
	},
}

var (
	tmplMarkdown  = texttmpl.Must(texttmpl.ParseFS(templates, "templ/markdown.tmpl")).Funcs(renderFns)
	tmplHTML      = htmltmpl.Must(htmltmpl.ParseFS(templates, "templ/html.tmpl")).Funcs(renderFns)
	tmplPlaintext = texttmpl.Must(texttmpl.ParseFS(templates, "templ/plaintext.tmpl")).Funcs(renderFns)
)

type template interface {
	Execute(wr io.Writer, data any) error
}

func templateRenderer(t template) func(renderContext, io.Writer) error {
	return func(c renderContext, out io.Writer) error {
		if err := t.Execute(out, c); err != nil {
			return fmt.Errorf("render template: %w", err)
		}
		return nil
	}
}
