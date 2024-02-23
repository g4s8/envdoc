package main

import (
	"fmt"
	"io"
)

type generator struct {
	fileName     string
	execLine     int
	targetType   string
	all          bool // generate documentation for all types in the file
	tmpl         template
	renderConfig renderConfig
	prefix       string
	noStyles     bool
	fieldNames   bool
}

type generatorOption func(*generator) error

func withType(targetType string) generatorOption {
	return func(g *generator) error {
		g.targetType = targetType
		g.all = false
		return nil
	}
}

func withAll() generatorOption {
	return func(g *generator) error {
		g.targetType = ""
		g.all = true
		return nil
	}
}

func withFormat(formatName string) generatorOption {
	return func(g *generator) error {
		switch formatName {
		case "":
			fallthrough
		case "markdown":
			g.tmpl = tmplMarkdown
			g.renderConfig = renderMarkdown
		case "plaintext":
			g.tmpl = tmplPlaintext
			g.renderConfig = renderPlaintext
		case "html":
			g.tmpl = tmplHTML
			g.renderConfig = renderHTML
		case "dotenv":
			g.tmpl = tmplDotEnv
			g.renderConfig = renderDotenv
		default:
			return fmt.Errorf("unknown format: %s", formatName)
		}
		return nil
	}
}

func withPrefix(prefix string) generatorOption {
	return func(g *generator) error {
		g.prefix = prefix
		return nil
	}
}

func withNoStyles() generatorOption {
	return func(g *generator) error {
		g.noStyles = true
		return nil
	}
}

func withFieldNames() generatorOption {
	return func(g *generator) error {
		g.fieldNames = true
		return nil
	}
}

func newGenerator(fileName string, execLine int, opts ...generatorOption) (*generator, error) {
	g := &generator{fileName: fileName, execLine: execLine}
	for _, opt := range opts {
		if err := opt(g); err != nil {
			return nil, err
		}
	}
	if g.tmpl == nil {
		g.tmpl = tmplMarkdown
		g.renderConfig = renderMarkdown
	}
	return g, nil
}

func (g *generator) generate(out io.Writer) error {
	insp := newInspector(g.targetType, g.all, g.execLine, g.fieldNames)
	data, err := insp.inspectFile(g.fileName)
	if err != nil {
		return fmt.Errorf("inspect file: %w", err)
	}
	renderer := templateRenderer(g.tmpl)
	rctx := newRenderContext(data, g.renderConfig, g.prefix, g.noStyles)
	if err := renderer(rctx, out); err != nil {
		return err
	}
	return nil
}
