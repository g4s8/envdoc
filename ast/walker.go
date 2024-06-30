package ast

import (
	"go/ast"
	"go/token"
)

func Walk(n ast.Node, fset *token.FileSet, h FileHandler) {
	v := newPkgVisitor(fset, h)
	ast.Walk(v, n)
}
