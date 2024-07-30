package ast

import "go/ast"

type fieldVisitor struct {
	pkg string
	h   FieldHandler

	nested bool
}

func newFieldVisitor(pkg string, h FieldHandler) *fieldVisitor {
	return &fieldVisitor{pkg: pkg, h: h}
}

func (v *fieldVisitor) Visit(n ast.Node) ast.Visitor {
	debugNode("field", n)
	switch t := n.(type) {
	case *ast.StructType:
		v.nested = true
		return v
	case *ast.Field:
		if !v.nested {
			return nil
		}
		fs := getFieldSpec(t, v.pkg)
		if fs == nil {
			return nil
		}
		if fa := v.h.onField(fs); fa != nil {
			return newFieldVisitor(v.pkg, fa)
		}
	}
	return v
}
