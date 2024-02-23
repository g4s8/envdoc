package main

import (
	"fmt"
	"io"
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

type renderItemConfigt struct {
	SeparatorFormat  string
	SeparatorDefault string
	OptRequired      string
	OptExpand        string
	OptNonEmpty      string
	OptFromFile      string
	EnvDefaultFormat string
}
type renderConfig struct {
	Item renderItemConfigt
}

type renderContext struct {
	Title    string
	Sections []renderSection
	Styles   bool
	Config   renderConfig
}

func newRenderContext(scopes []*EnvScope, cfg renderConfig, envPrefix string, noStyles bool) renderContext {
	res := renderContext{
		Sections: make([]renderSection, len(scopes)),
		Styles:   !noStyles,
		Config:   cfg,
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

func newRenderItem(item *EnvDocItem, envPrefix string) renderItem {
	log := logger()
	children := make([]renderItem, len(item.Children))
	log.Printf("render item %s", item.Name)
	for i, child := range item.Children {
		log.Printf("render child item %s", child.Name)
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

var (
	tmplMarkdown  = newTmplText("markdown.tmpl")
	tmplHTML      = newTmplText("html.tmpl")
	tmplPlaintext = newTmplText("plaintext.tmpl")
	tmplDotEnv    = newTmplText("dotenv.tmpl")

	renderMarkdown = renderConfig{
		Item: renderItemConfigt{
			SeparatorFormat:  "separated by `%s`",
			SeparatorDefault: "comma-separated",
			OptRequired:      "**required**",
			OptExpand:        "expand",
			OptFromFile:      "from-file",
			OptNonEmpty:      "non-empty",
			EnvDefaultFormat: "default: `%s`",
		},
	}
	renderPlaintext = renderConfig{
		Item: renderItemConfigt{
			SeparatorFormat:  "separated by `%s`",
			SeparatorDefault: "comma-separated",
			OptRequired:      "required",
			OptExpand:        "expand",
			OptFromFile:      "from-file",
			OptNonEmpty:      "non-empty",
			EnvDefaultFormat: "default: `%s`",
		},
	}
	renderDotenv = renderConfig{
		Item: renderItemConfigt{
			SeparatorFormat:  "separated by '%s'",
			SeparatorDefault: "comma-separated",
			OptRequired:      "required",
			OptExpand:        "expand",
			OptFromFile:      "from-file",
			OptNonEmpty:      "non-empty",
			EnvDefaultFormat: "default: '%s'",
		},
	}
	renderHTML = renderConfig{
		Item: renderItemConfigt{
			SeparatorFormat:  `separated by "<code>%s</code>"`,
			SeparatorDefault: "comma-separated",
			OptRequired:      "<strong>required</strong>",
			OptExpand:        "expand",
			OptFromFile:      "from-file",
			OptNonEmpty:      "non-empty",
			EnvDefaultFormat: "default: <code>%s</code>",
		},
	}
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
