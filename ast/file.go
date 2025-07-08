package ast

import (
	"go/ast"
	"go/doc"
	"go/token"
	"strings"

	"github.com/g4s8/envdoc/debug"
)

type fileVisitorHandler = interface {
	TypeHandler
	CommentHandler
	ImportHandler
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
	case *ast.ImportSpec:
		var spec ImportSpec
		if n := t.Name; n != nil {
			spec.Name = n.Name
		}
		path := strings.TrimPrefix(t.Path.Value, "\"")
		path = strings.TrimSuffix(path, "\"")
		spec.Path = path
		debug.Logf("# V: import %q, name=%q\n", spec.Path, spec.Name)
		v.h.addImport(&spec)
		return nil
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
