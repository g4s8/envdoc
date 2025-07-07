package resolver

import (
	"fmt"
	"io"

	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/debug"
)

type typeQualifier struct {
	pkg  string
	name string
}

func (q typeQualifier) String() string {
	return fmt.Sprintf("%s.%s", q.pkg, q.name)
}

type TypeResolver struct {
	types map[typeQualifier]*ast.TypeSpec
}

func NewTypeResolver() *TypeResolver {
	return &TypeResolver{
		types: make(map[typeQualifier]*ast.TypeSpec),
	}
}

func (r *TypeResolver) AddTypes(pkg string, types []*ast.TypeSpec) {
	for _, t := range types {
		r.types[typeQualifier{pkg: pkg, name: t.Name}] = t
	}
}

func (r *TypeResolver) Resolve(f *ast.FileSpec, ref *ast.FieldTypeRef) *ast.TypeSpec {
	pkg := ref.Pkg
	if pkg != "" {
		for _, alias := range f.Imports {
			if alias.Name == pkg {
				pkg = alias.PathName()
				break
			}
		}
	}
	tq := typeQualifier{pkg: pkg, name: ref.Name}
	ts := r.types[tq]
	debug.Logf("# RES: ref=%q tq=%q ts=%q",
		ref, tq, ts)
	return ts
}

func ResolveAllTypes(files []*ast.FileSpec) *TypeResolver {
	r := NewTypeResolver()
	for _, f := range files {
		pkg := f.Pkg
		r.AddTypes(pkg, f.Types)
	}
	return r
}

func (r *TypeResolver) Debug(out io.Writer) {
	fmt.Fprintln(out, "Resolved types:")
	for k, v := range r.types {
		fmt.Fprintf(out, "  %s.%s: %q (export=%t)\n",
			k.pkg, k.name, v.Name, v.Export)
	}
}
