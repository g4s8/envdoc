package main

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestASTTypeDocResolver(t *testing.T) {
	t.Run("Fail", func(t *testing.T) {
		fset := token.NewFileSet()
		astFile := ast.File{}
		_, err := newASTTypeDocResolver(fset, &astFile)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		t.Logf("Error: %v", err)
	})
}
