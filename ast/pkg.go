package ast

import (
	"go/ast"
	"go/doc"
	"go/token"
)

type pkgVisitor struct {
	fset *token.FileSet
	h    FileHandler

	pkg  string
	docs *doc.Package
}

func newPkgVisitor(fset *token.FileSet, h FileHandler) *pkgVisitor {
	return &pkgVisitor{
		fset: fset,
		h:    h,
	}
}

func (p *pkgVisitor) Visit(n ast.Node) ast.Visitor {
	debugNode("pkg", n)
	switch t := n.(type) {
	//nolint:staticcheck
	case *ast.Package:
		p.pkg = t.Name
		p.docs = doc.New(t, "./", doc.PreserveAST|doc.AllDecls)
		return p
	case *ast.File:
		f := p.fset.File(t.Pos())
		if fa := p.h.onFile(&FileSpec{
			Name: f.Name(),
			Pkg:  p.pkg,
		}); fa != nil {
			return newFileVisitor(p.fset, t, p.docs, fa)
		}
	}
	return p
}
