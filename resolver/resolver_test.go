package resolver

import (
	"testing"

	"github.com/g4s8/envdoc/ast"
)

func TestResolver(t *testing.T) { //nolint:cyclop
	res := ResolveAllTypes([]*ast.FileSpec{
		{
			Pkg: "main",
			Types: []*ast.TypeSpec{
				{
					Name:   "Foo",
					Export: true,
				},
				{
					Name:   "Bar",
					Export: false,
				},
			},
		},
		{
			Pkg: "test",
			Types: []*ast.TypeSpec{
				{
					Name:   "Baz",
					Export: true,
				},
			},
		},
	})
	foo := res.Resolve(&ast.FileSpec{}, &ast.FieldTypeRef{Pkg: "main", Name: "Foo"})
	if foo == nil {
		t.Fatalf("Foo type not resolved")
	}
	if foo.Name != "Foo" {
		t.Errorf("Invalid Foo type: %s", foo.Name)
	}

	bar := res.Resolve(&ast.FileSpec{}, &ast.FieldTypeRef{Pkg: "main", Name: "Bar"})
	if bar == nil {
		t.Fatalf("Bar type not resolved")
	}
	if bar != nil && bar.Name != "Bar" {
		t.Errorf("Invalid Bar type: %s", bar.Name)
	}

	baz := res.Resolve(&ast.FileSpec{}, &ast.FieldTypeRef{Pkg: "test", Name: "Baz"})
	if baz == nil {
		t.Fatalf("Baz type not resolved")
	}
	if baz.Name != "Baz" {
		t.Errorf("Invalid Baz type: %s", baz.Name)
	}

	bazAlias := res.Resolve(&ast.FileSpec{
		Imports: []*ast.ImportSpec{
			{
				Name: "tst",
				Path: "example.com/foo/test",
			},
		},
	}, &ast.FieldTypeRef{Pkg: "tst", Name: "Baz"})
	if bazAlias == nil {
		t.Fatalf("Baz type with alias not resolved")
	}
	if bazAlias.Name != "Baz" {
		t.Errorf("Invalid Baz type with alias: %s", bazAlias.Name)
	}

	nope := res.Resolve(&ast.FileSpec{}, &ast.FieldTypeRef{Pkg: "test", Name: "Nope"})
	if nope != nil {
		t.Errorf("Nope type resolved, but it should not")
	}

	wrongPgk := res.Resolve(&ast.FileSpec{}, &ast.FieldTypeRef{Pkg: "main", Name: "Baz"})
	if wrongPgk != nil {
		t.Errorf("Baz type resolved, but it should not")
	}
}
