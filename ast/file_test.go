package ast

import (
	"go/ast"
	"path/filepath"
	"testing"
)

func TestFileVisitor(t *testing.T) {
	fset, pkg, docs := loadTestFileSet(t)
	// parser.ParseDir keys files with OS-native separators, so build the
	// lookup path accordingly (backslash on Windows).
	fh, fv, file := testFileVisitor(fset, pkg, filepath.Join("testdata", "onetype.go"), docs)
	ast.Walk(fv, file)

	types := make([]*TypeSpec, 0)
	for _, f := range fh.files {
		types = append(types, f.Types...)
	}
	if expect, actual := 1, len(types); expect != actual {
		t.Fatalf("expected %d types, got %d", expect, actual)
	}
	if expect, actual := "One", types[0].Name; expect != actual {
		t.Fatalf("expected type name %q, got %q", expect, actual)
	}
}
