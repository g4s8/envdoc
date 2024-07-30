//go:build !coverage

package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/debug"
)

var DebugConfig struct {
	Enabled bool
}

func printTraverse(files []*ast.FileSpec, level int) {
	indent := strings.Repeat("  ", level)
	for _, file := range files {
		fmt.Printf("%sFILE:%q\n", indent, file.Name)
		printTraverseTypes(file.Types, level+1)
	}
}

func printTraverseTypes(types []*ast.TypeSpec, level int) {
	indent := strings.Repeat("  ", level)
	for _, t := range types {
		fmt.Printf("%sTYPE:%q; doc: %q\n", indent, t.Name, t.Doc)
		printTraverseFields(t.Fields, level+1)
	}
}

func printTraverseFields(fields []*ast.FieldSpec, level int) {
	indent := strings.Repeat("  ", level)
	for _, f := range fields {
		names := strings.Join(f.Names, ", ")
		fmt.Printf("%sFIELD:%s (%s); doc: %q\n", indent, names, f.TypeRef.String(), f.Doc)
		printTraverseFields(f.Fields, level+1)
	}
}

func (r *TypeResolver) fprint(out io.Writer) {
	fmt.Fprintln(out, "Resolved types:")
	for k, v := range r.types {
		fmt.Fprintf(out, "  %s.%s: %q (export=%t)\n",
			k.pkg, k.name, v.Name, v.Export)
	}
}

func printScopesTree(s []*EnvScope) {
	if !debug.Config.Enabled {
		return
	}
	debug.Log("Scopes tree:\n")
	for _, scope := range s {
		debug.Logf(" - %q\n", scope.Name)
		for _, item := range scope.Vars {
			printDocItem("  ", item)
		}
	}
}

func printDocItem(prefix string, item *EnvDocItem) {
	debug.Logf("%s- %q\n", prefix, item.Name)
	for _, child := range item.Children {
		printDocItem(prefix+"  ", child)
	}
}
