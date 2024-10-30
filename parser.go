package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"

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

	col := ast.NewRootCollector(p.dir, colOpts...)

	if p.debug {
		fmt.Printf("Parsing dir %q (f=%q t=%q)\n", p.dir, p.fileGlob, p.typeGlob)
	}
	// walk through the directory and each subdirectory and call parseDir for each of them
	if err := filepath.Walk(p.dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk through dir: %w", err)
		}
		if !info.IsDir() {
			return nil
		}
		if err := parseDir(path, fset, matcher, col); err != nil {
			return fmt.Errorf("failed to parse dir %q: %w", path, err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk through dir: %w", err)
	}

	if p.debug {
		fmt.Printf("Resolved types:\n")
		printTraverse(col.Files(), 0)
	}

	return col.Files(), nil
}

func parseDir(dir string, fset *token.FileSet, matcher func(fs.FileInfo) bool, col *ast.RootCollector) error {
	pkgs, err := parser.ParseDir(fset, dir, matcher, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return fmt.Errorf("failed to parse dir: %w", err)
	}

	for _, pkg := range pkgs {
		ast.Walk(pkg, fset, col)
	}
	return nil
}
