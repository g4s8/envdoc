package render

import (
	"fmt"
	"io"

	"github.com/g4s8/envdoc/types"
)

type Renderer struct {
	format   types.OutFormat
	noStyles bool
}

func NewRenderer(format types.OutFormat, noStyles bool) *Renderer {
	return &Renderer{
		format:   format,
		noStyles: noStyles,
	}
}

func (r *Renderer) Render(scopes []*types.EnvScope, out io.Writer) error {
	cfg, ok := configs[r.format]
	if !ok {
		return fmt.Errorf("unknown format: %q", r.format)
	}

	c := newRenderContext(scopes, cfg, r.noStyles)
	f := templateRenderer(cfg.tmpl)

	if err := f(c, out); err != nil {
		return fmt.Errorf("render: %w", err)
	}
	return nil
}

type renderSection struct {
	Name  string
	Doc   string
	Items []renderItem
}

type renderItem struct {
	EnvName      string `json:"env_name"`
	Doc          string `json:"doc"`
	EnvDefault   string `json:"env_default,omitempty"`
	EnvSeparator string `json:"env_separator,omitempty"`

	Required bool `json:"required,omitempty"`
	Expand   bool `json:"expand,omitempty"`
	NonEmpty bool `json:"non_empty,omitempty"`
	FromFile bool `json:"from_file,omitempty"`

	children []renderItem
	Indent   int `json:"-"`
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
	Config   renderConfig
}

func newRenderContext(scopes []*types.EnvScope, cfg renderConfig, noStyles bool) renderContext {
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
			item := newRenderItem(item)
			item.Indent = 1
			section.Items[j] = item
		}
		res.Sections[i] = section
	}
	return res
}

func newRenderItem(item *types.EnvDocItem) renderItem {
	children := make([]renderItem, len(item.Children))
	for i, child := range item.Children {
		children[i] = newRenderItem(child)
	}
	return renderItem{
		EnvName:      item.Name,
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
