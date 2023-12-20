package main

import (
	"fmt"
	"io"
)

type generator struct {
	fileName   string
	execLine   int
	targetType string
	tmpl       template
}

type generatorOption func(*generator) error

func withType(targetType string) generatorOption {
	return func(g *generator) error {
		g.targetType = targetType
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
		case "plaintext":
			g.tmpl = tmplPlaintext
		case "html":
			g.tmpl = tmplHTML
		default:
			return fmt.Errorf("unknown format: %s", formatName)
		}
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
		return nil, fmt.Errorf("format is not specified")
	}
	return g, nil
}

func (g *generator) generate(out io.Writer) error {
	insp := newInspector(g.targetType, g.execLine)
	err, data := insp.inspectFile(g.fileName)
	if err != nil {
		return fmt.Errorf("inspect file: %w", err)
	}
	renderer := templateRenderer(g.tmpl)
	rctx := newRenderContext(data)
	if err := renderer(rctx, out); err != nil {
		return err
	}
	return nil
}
