package ast

import (
	"go/ast"
)

type typeVisitorHandler = interface {
	CommentHandler
	FieldHandler
}

type typeVisitor struct {
	pkg string
	h   typeVisitorHandler
}

func newTypeVisitor(pkg string, h typeVisitorHandler) *typeVisitor {
	return &typeVisitor{pkg: pkg, h: h}
}

func (v *typeVisitor) Visit(n ast.Node) ast.Visitor {
	debugNode("type", n)
	switch t := n.(type) {
	case *ast.Comment:
		v.h.setComment(&CommentSpec{
			Text: t.Text,
		})
		return nil
	case *ast.Field:
		fs := getFieldSpec(t, v.pkg)
		if fs == nil {
			return nil
		}
		if fa := v.h.onField(fs); fa != nil {
			return newFieldVisitor(v.pkg, fa)
		}
		return nil
	}
	return v
}
