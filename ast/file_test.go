package ast

import (
	"go/ast"
	"testing"
)

type testFileVisitorHandler struct {
	comments []*CommentSpec
	types    []*TypeSpec
}

func (h *testFileVisitorHandler) setComment(c *CommentSpec) {
	h.comments = append(h.comments, c)
}

func (h *testFileVisitorHandler) onType(t *TypeSpec) typeVisitorHandler {
	h.types = append(h.types, t)
	return nil
}

func TestFileVisitor(t *testing.T) {
	fset, pkg, docs := loadTestFileSet(t)
	h := &testFileVisitorHandler{}
	file := pkg.Files["testdata/onetype.go"]
	v := newFileVisitor(fset, file, docs, h)
	ast.Walk(v, file)

	if expect, actual := 1, len(h.comments); expect != actual {
		t.Fatalf("expected %d comments, got %d", expect, actual)
	}
	checkCommentsEq(t, &CommentSpec{
		Text: "onetype",
	}, h.comments[0])
	if expect, actual := 1, len(h.types); expect != actual {
		t.Fatalf("expected %d types, got %d", expect, actual)
	}
	if expect, actual := "One", h.types[0].Name; expect != actual {
		t.Fatalf("expected type name %q, got %q", expect, actual)
	}
}
