package main

import (
	"fmt"
	"io"
	"os"

	"github.com/g4s8/envdoc/ast"
)

type Generator struct {
	parser    *ast.Parser
	converter *Converter
	renderer  *Renderer
}

func NewGenerator(parser *ast.Parser, converter *Converter, renderer *Renderer) *Generator {
	return &Generator{
		parser:    parser,
		converter: converter,
		renderer:  renderer,
	}
}

func (g *Generator) Generate(dir string, out io.Writer) error {
	files, err := g.parser.Parse(dir)
	if err != nil {
		return fmt.Errorf("parse dir: %w", err)
	}

	res := ResolveAllTypes(files)
	if DebugConfig.Enabled {
		res.fprint(os.Stdout)
	}

	scopes := g.converter.ScopesFromFiles(res, files)
	if DebugConfig.Enabled {
		printScopesTree(scopes)
	}

	if err := g.renderer.Render(scopes, out); err != nil {
		return fmt.Errorf("render: %w", err)
	}

	return nil
}
