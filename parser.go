package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"

	"github.com/g4s8/envdoc/ast"
)

type parserConfigOption func(*Parser)

func withDebug(debug bool) parserConfigOption {
	return func(p *Parser) {
		p.debug = debug
	}
}

func withExecConfig(execFile string, execLine int) parserConfigOption {
	return func(p *Parser) {
		p.gogenFile = execFile
		p.gogenLine = execLine
	}
}

type Parser struct {
	dir       string
	fileGlob  string
	typeGlob  string
	gogenLine int
	gogenFile string
	debug     bool
}

func NewParser(dir, fileGlob, typeGlob string, opts ...parserConfigOption) *Parser {
	p := &Parser{
		dir:      dir,
		fileGlob: fileGlob,
		typeGlob: typeGlob,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Parser) Parse() ([]*ast.FileSpec, error) {
	fset := token.NewFileSet()
	var matcher func(fs.FileInfo) bool
	if p.fileGlob != "" {
		m, err := newGlobFileMatcher(p.fileGlob)
		if err != nil {
			return nil, err
		}
		matcher = m
	}

	pkgs, err := parser.ParseDir(fset, p.dir, matcher, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dir: %w", err)
	}

	var colOpts []ast.RootCollectorOption
	if p.typeGlob == "" {
		colOpts = append(colOpts, ast.WithGoGenDecl(p.gogenLine, p.gogenFile))
	} else {
		m, err := newGlobMatcher(p.typeGlob)
		if err != nil {
			return nil, fmt.Errorf("create type glob matcher: %w", err)
		}
		colOpts = append(colOpts, ast.WithTypeGlob(m))
	}
	if p.fileGlob != "" {
		m, err := newGlobMatcher(p.fileGlob)
		if err != nil {
			return nil, fmt.Errorf("create file glob matcher: %w", err)
		}
		colOpts = append(colOpts, ast.WithFileGlob(m))
	}

	col := ast.NewRootCollector(colOpts...)
	for _, pkg := range pkgs {
		ast.Walk(pkg, fset, col)
	}

	if p.debug {
		printTraverse(col.Files(), 0)
	}

	return col.Files(), nil
}
