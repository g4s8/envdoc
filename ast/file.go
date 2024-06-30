package ast

import (
	"go/ast"
	"go/doc"
	"go/token"
	"strings"
)

type fileVisitorHandler = interface {
	TypeHandler
	CommentHandler
}

type fileVisitor struct {
	fset *token.FileSet
	file *ast.File
	docs *doc.Package
	h    fileVisitorHandler
}

func newFileVisitor(fset *token.FileSet, file *ast.File, docs *doc.Package, h fileVisitorHandler) *fileVisitor {
	return &fileVisitor{
		fset: fset,
		file: file,
		docs: docs,
		h:    h,
	}
}

func (v *fileVisitor) Visit(n ast.Node) ast.Visitor {
	debugNode("file", n)
	switch t := n.(type) {
	case *ast.Comment:
		line := findCommentLine(t, v.fset, v.file)
		text := strings.TrimPrefix(t.Text, "//")
		text = strings.TrimSpace(text)
		v.h.setComment(&CommentSpec{
			Line: line,
			Text: text,
		})
		return nil
	case *ast.TypeSpec:
		doc := resolveTypeDocs(v.docs, t)
		if ta := v.h.onType(&TypeSpec{
			Name: t.Name.Name,
			Doc:  doc,
		}); ta != nil {
			return newTypeVisitor(v.file.Name.String(), ta)
		}
		return nil
	}
	return v
}
