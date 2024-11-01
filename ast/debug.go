//go:build !coverage

package ast

import (
	"fmt"
	"strings"
)

func printTraverse(files []*FileSpec, level int) {
	indent := strings.Repeat("  ", level)
	for _, file := range files {
		fmt.Printf("%sFILE:%q\n", indent, file.Name)
		printTraverseTypes(file.Types, level+1)
	}
}

func printTraverseTypes(types []*TypeSpec, level int) {
	indent := strings.Repeat("  ", level)
	for _, t := range types {
		fmt.Printf("%sTYPE:%q; doc: %q\n", indent, t.Name, t.Doc)
		printTraverseFields(t.Fields, level+1)
	}
}

func printTraverseFields(fields []*FieldSpec, level int) {
	indent := strings.Repeat("  ", level)
	for _, f := range fields {
		names := strings.Join(f.Names, ", ")
		fmt.Printf("%sFIELD:%s (%s); doc: %q\n", indent, names, f.TypeRef.String(), f.Doc)
		printTraverseFields(f.Fields, level+1)
	}
}
