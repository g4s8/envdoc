package ast

import (
	"go/ast"
	"slices"
	"testing"
)

func TestPkgVisitor(t *testing.T) {
	fset, pkg, _ := loadTestFileSet(t)
	h := &testFileHandler{}
	v := newPkgVisitor(fset, h)
	ast.Walk(v, pkg)
	if len(h.files) != 4 {
		t.Fatalf("expected 4 files, got %d", len(h.files))
	}
	expectFiles := []string{"empty.go", "onetype.go", "twotypes.go", "fields.go"}
	expectPkg := "testdata"
	fileNames := make([]string, len(h.files))
	for i, f := range h.files {
		fileNames[i] = f.Name
		t.Logf("file %q", f.Name)
	}
	for _, e := range expectFiles {
		e = "testdata/" + e
		if !slices.Contains(fileNames, e) {
			t.Fatalf("file %q not found", e)
		}
	}
	for _, f := range h.files {
		if f.Pkg != expectPkg {
			t.Fatalf("expected pkg %q, got %q", expectPkg, f.Pkg)
		}
	}
}
