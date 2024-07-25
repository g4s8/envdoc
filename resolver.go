package main

import (
	"github.com/g4s8/envdoc/ast"
)

type typeQualifier struct {
	pkg  string
	name string
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

func (r *TypeResolver) Resolve(ref *ast.FieldTypeRef) *ast.TypeSpec {
	return r.types[typeQualifier{pkg: ref.Pkg, name: ref.Name}]
}

func ResolveAllTypes(files []*ast.FileSpec) *TypeResolver {
	r := NewTypeResolver()
	for _, f := range files {
		pkg := f.Pkg
		r.AddTypes(pkg, f.Types)
	}
	return r
}
