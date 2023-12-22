package main

import (
	"embed"
	"fmt"
	"io"

	htmltmpl "html/template"
	texttmpl "text/template"
)

type renderSection struct {
	Name  string
	Doc   string
	Items []renderItem
}

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
	Title    string
	Sections []renderSection
}

func newRenderContext(scopes []*EnvScope, envPrefix string) renderContext {
	res := renderContext{
		Sections: make([]renderSection, len(scopes)),
	}
	res.Title = "Environment Variables"
	for i, scope := range scopes {
		section := renderSection{
			Name:  scope.Name,
			Doc:   scope.Doc,
			Items: make([]renderItem, len(scope.Vars)),
		}
		for j, item := range scope.Vars {
			section.Items[j] = renderItem{
				EnvName:      fmt.Sprintf("%s%s", envPrefix, item.Name),
				Doc:          item.Doc,
				EnvDefault:   item.Opts.Default,
				EnvSeparator: item.Opts.Separator,
				Required:     item.Opts.Required,
				Expand:       item.Opts.Expand,
				NonEmpty:     item.Opts.NonEmpty,
				FromFile:     item.Opts.FromFile,
			}
		}
		res.Sections[i] = section
	}
	return res
}

//go:embed templ
var templates embed.FS

var (
	tmplMarkdown  = texttmpl.Must(texttmpl.ParseFS(templates, "templ/markdown.tmpl"))
	tmplHTML      = htmltmpl.Must(htmltmpl.ParseFS(templates, "templ/html.tmpl"))
	tmplPlaintext = texttmpl.Must(texttmpl.ParseFS(templates, "templ/plaintext.tmpl"))
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
