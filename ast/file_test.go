package ast

import (
	"go/ast"
	"testing"
)

func TestFileVisitor(t *testing.T) {
	fset, pkg, docs := loadTestFileSet(t)
	fh, fv, file := testFileVisitor(fset, pkg, "testdata/onetype.go", docs)
	ast.Walk(fv, file)

	if expect, actual := 1, len(fh.comments); expect != actual {
		t.Fatalf("expected %d comments, got %d", expect, actual)
	}
	checkCommentsEq(t, &CommentSpec{
		Text: "onetype",
	}, fh.comments[0])
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
