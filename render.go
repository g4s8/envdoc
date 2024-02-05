package main

import (
	"embed"
	"fmt"
	"io"
	"strings"

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

	children []renderItem
	Indent   int
}

func (i renderItem) Children(indentInc int) []renderItem {
	indent := i.Indent + indentInc
	res := make([]renderItem, len(i.children))
	for j, child := range i.children {
		child.Indent = indent
		res[j] = child
	}
	return res
}

type renderContext struct {
	Title    string
	Sections []renderSection
	Styles   bool
}

func newRenderContext(scopes []*EnvScope, envPrefix string, noStyles bool) renderContext {
	res := renderContext{
		Sections: make([]renderSection, len(scopes)),
		Styles:   !noStyles,
	}
	res.Title = "Environment Variables"
	for i, scope := range scopes {
		section := renderSection{
			Name:  scope.Name,
			Doc:   scope.Doc,
			Items: make([]renderItem, len(scope.Vars)),
		}
		for j, item := range scope.Vars {
			item := newRenderItem(item, envPrefix)
			item.Indent = 1
			section.Items[j] = item
		}
		res.Sections[i] = section
	}
	return res
}

func newRenderItem(item EnvDocItem, envPrefix string) renderItem {
	children := make([]renderItem, len(item.Children))
	debug("render item %s", item.Name)
	for i, child := range item.Children {
		debug("render child item %s", child.Name)
		children[i] = newRenderItem(child, envPrefix)
	}
	return renderItem{
		EnvName:      fmt.Sprintf("%s%s", envPrefix, item.Name),
		Doc:          item.Doc,
		EnvDefault:   item.Opts.Default,
		EnvSeparator: item.Opts.Separator,
		Required:     item.Opts.Required,
		Expand:       item.Opts.Expand,
		NonEmpty:     item.Opts.NonEmpty,
		FromFile:     item.Opts.FromFile,
		children:     children,
	}
}

//go:embed templ
var templates embed.FS

var tplFuncs = map[string]any{
	"repeat": strings.Repeat,
}

func _() {
	// texttmpl.ParseFS
}

var (
	tmplMarkdown  = texttmpl.Must(texttmpl.New("markdown.tmpl").Funcs(tplFuncs).ParseFS(templates, "templ/markdown.tmpl"))
	tmplHTML      = htmltmpl.Must(htmltmpl.New("html.tmpl").Funcs(tplFuncs).ParseFS(templates, "templ/html.tmpl"))
	tmplPlaintext = texttmpl.Must(texttmpl.New("plaintext.tmpl").Funcs(tplFuncs).ParseFS(templates, "templ/plaintext.tmpl"))
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
