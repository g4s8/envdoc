package ast

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
)

type T interface {
	Helper()
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

func loadTestFileSet(t T) (*token.FileSet, *ast.Package, *doc.Package) {
	t.Helper()
	// load go files from ./testdata dir
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "./testdata", nil, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("failed to parse testdata: %s", err)
	}
	pkg, ok := pkgs["testdata"]
	if !ok {
		t.Fatal("package 'testdata' not found")
	}
	docs := doc.New(pkg, "./", doc.PreserveAST|doc.AllDecls)
	return fset, pkg, docs
}

func checkCommentsEq(t T, expect, actual *CommentSpec) {
	t.Helper()

	if expect.Text != actual.Text {
		t.Fatalf("expected comment %q, got %q", expect.Text, actual.Text)
	}
	if expect.Line != actual.Line {
		t.Fatalf("expected comment line %d, got %d", expect.Line, actual.Line)
	}
}
