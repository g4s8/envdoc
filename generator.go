package main

import (
	"fmt"
	"io"
)

type generator struct {
	fileName   string
	execLine   int
	targetType string
	format     itemFormat
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
			g.format = fmtMD
		case "plaintext":
			g.format = fmtPlain
		case "html":
			g.format = fmtHTML
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
	if g.format == nil {
		return nil, fmt.Errorf("format is not specified")
	}
	return g, nil
}

func (g *generator) generate(out io.Writer) error {
	output := newDocOutput(out, g.format)
	output.begin()
	insp := newInspector(g.targetType, output, g.execLine)
	if err := insp.inspectFile(g.fileName); err != nil {
		return fmt.Errorf("inspect file: %w", err)
	}
	output.end()
	return nil
}
