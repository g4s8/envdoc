package ast

import (
	"go/ast"
	"testing"
)

type fileHandler struct {
	types []*TypeSpec
	typeH *typeHandler
}

func (h *fileHandler) setComment(_ *CommentSpec) {
}

func (h *fileHandler) onType(t *TypeSpec) typeVisitorHandler {
	h.types = append(h.types, t)
	h.typeH = &typeHandler{}
	return h.typeH
}

type typeHandler struct {
	fields   []*FieldSpec
	comments []*CommentSpec
}

func (h *typeHandler) setComment(c *CommentSpec) {
	h.comments = append(h.comments, c)
}

func (h *typeHandler) onField(f *FieldSpec) FieldHandler {
	h.fields = append(h.fields, f)
	return nil
}

func TestTypesVisitor(t *testing.T) {
	fset, pkg, docs := loadTestFileSet(t)
	file := pkg.Files["testdata/fields.go"]
	h := &fileHandler{}
	v := newFileVisitor(fset, file, docs, h)

	ast.Walk(v, file)

	fh := h.typeH
	if expect, actual := 3, len(fh.fields); expect != actual {
		t.Fatalf("expected %d fields, got %d", expect, actual)
	}
	if expect, actual := "A", fh.fields[0].Names[0]; expect != actual {
		t.Fatalf("expected field name %q, got %q", expect, actual)
	}
	if expect, actual := "B", fh.fields[1].Names[0]; expect != actual {
		t.Fatalf("expected field name %q, got %q", expect, actual)
	}
	if expect, actual := "C", fh.fields[2].Names[0]; expect != actual {
		t.Fatalf("expected field name %q, got %q", expect, actual)
	}
}
