package main

import (
	"fmt"
	"testing"

	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/resolver"
	"github.com/g4s8/envdoc/types"
)

var opts = ConverterOpts{
	TagName:    "env",
	TagDefault: "envDefault",
}

func TestConvertDocItems(t *testing.T) {
	opts := opts
	opts.UseFieldNames = true

	c := NewConverter(types.TargetTypeCaarlos0, opts)
	fieldValues := []*ast.FieldSpec{
		{
			Names: []string{"Field1"},
			TypeRef: ast.FieldTypeRef{
				Name: "string",
				Kind: ast.FieldTypeIdent,
			},
			Tag: `env:"FIELD1,required,file"`,
			Doc: "Field1 doc",
		},
		{
			Names: []string{"Field2", "Field3"},
			TypeRef: ast.FieldTypeRef{
				Name: "int",
				Kind: ast.FieldTypeIdent,
			},
			Doc: "Field2 and Field3 doc",
		},
		{
			Names: []string{"FieldDef"},
			TypeRef: ast.FieldTypeRef{
				Name: "string",
				Kind: ast.FieldTypeIdent,
			},
			Doc: "Field with default",
			Tag: `env:"FIELD_DEF" envDefault:"envdef"`,
		},
		{
			Names: []string{"FieldArr"},
			TypeRef: ast.FieldTypeRef{
				Name: "[]string",
				Kind: ast.FieldTypeArray,
			},
			Doc: "Field array",
			Tag: `env:"FIELD_ARR"`,
		},
		{
			Names: []string{"FieldArrSep"},
			TypeRef: ast.FieldTypeRef{
				Name: "[]string",
				Kind: ast.FieldTypeArray,
			},
			Doc: "Field array with separator",
			Tag: `env:"FIELD_ARR_SEP" envSeparator:":"`,
		},
		{
			Names: []string{"FooField"},
			TypeRef: ast.FieldTypeRef{
				Name: "Foo",
				Kind: ast.FieldTypePtr,
			},
			Tag: `envPrefix:"FOO_"`,
		},
		{
			Names: []string{"BarField"},
			TypeRef: ast.FieldTypeRef{
				Pkg:  "config",
				Name: "Bar",
				Kind: ast.FieldTypeIdent,
			},
			Tag: `envPrefix:"BAR_"`,
		},
		{
			Names: []string{"StructField"},
			TypeRef: ast.FieldTypeRef{
				Kind: ast.FieldTypeStruct,
			},
			Fields: []*ast.FieldSpec{
				{
					Names: []string{"Field1"},
					TypeRef: ast.FieldTypeRef{
						Name: "string",
						Kind: ast.FieldTypeIdent,
					},
					Doc: "Field1 doc",
					Tag: `env:"FIELD1"`,
				},
			},
			Tag: `envPrefix:"STRUCT_"`,
		},
		{
			Names: []string{},
			Doc:   "Embedded field",
			Fields: []*ast.FieldSpec{
				{
					Names: []string{"Field4"},
					TypeRef: ast.FieldTypeRef{
						Name: "string",
						Kind: ast.FieldTypeIdent,
					},
					Doc: "Field4 doc",
					Tag: `env:"FIELD4,notEmpty,expand"`,
				},
			},
			TypeRef: ast.FieldTypeRef{
				Kind: ast.FieldTypeStruct,
			},
		},
	}
	resolver := resolver.NewTypeResolver()
	resolver.AddTypes("", []*ast.TypeSpec{
		{
			Name: "Foo",
			Doc:  "Foo doc",
			Fields: []*ast.FieldSpec{
				{
					Names: []string{"FOne"},
					Doc:   "Foo one field",
					Tag:   `env:"F1"`,
				},
			},
		},
	})
	resolver.AddTypes("config", []*ast.TypeSpec{
		{
			Name: "Bar",
			Doc:  "Bar doc",
			Fields: []*ast.FieldSpec{
				{
					Names: []string{"BOne"},
					Doc:   "Bar one field",
					Tag:   `env:"B1"`,
				},
			},
		},
	})

	res := c.DocItemsFromFields(resolver, "", fieldValues)
	expect := []*types.EnvDocItem{
		{
			Name: "FIELD1",
			Doc:  "Field1 doc",
			Opts: types.EnvVarOptions{
				Required: true,
				FromFile: true,
			},
		},
		{
			Name: "FIELD2",
			Doc:  "Field2 and Field3 doc",
		},
		{
			Name: "FIELD3",
			Doc:  "Field2 and Field3 doc",
		},
		{
			Name: "FIELD_DEF",
			Doc:  "Field with default",
			Opts: types.EnvVarOptions{
				Default: "envdef",
			},
		},
		{
			Name: "FIELD_ARR",
			Doc:  "Field array",
			Opts: types.EnvVarOptions{
				Separator: ",",
			},
		},
		{
			Name: "FIELD_ARR_SEP",
			Doc:  "Field array with separator",
			Opts: types.EnvVarOptions{
				Separator: ":",
			},
		},
		{
			Name: "FOO_FIELD",
			Children: []*types.EnvDocItem{
				{
					Name: "FOO_F1",
					Doc:  "Foo one field",
				},
			},
		},
		{
			Name: "BAR_FIELD",
			Children: []*types.EnvDocItem{
				{
					Name: "BAR_B1",
					Doc:  "Bar one field",
				},
			},
		},
		{
			Name: "STRUCT_FIELD",
			Children: []*types.EnvDocItem{
				{
					Name: "STRUCT_FIELD1",
					Doc:  "Field1 doc",
				},
			},
		},
		{
			Name: "FIELD4",
			Doc:  "Field4 doc",
			Opts: types.EnvVarOptions{
				Required: true,
				NonEmpty: true,
				Expand:   true,
			},
		},
	}
	if len(expect) != len(res) {
		t.Errorf("Expected %d items, got %d", len(expect), len(res))
		for i, item := range expect {
			t.Logf("Expect[%d] %q", i, item.Name)
		}
		for i, item := range res {
			t.Logf("Actual[%d] %q", i, item.Name)
		}
		t.FailNow()
	}
	for i, item := range expect {
		checkDocItem(t, fmt.Sprintf("%d", i), item, res[i])
	}
}

func TestConverterScopes(t *testing.T) {
	files := []*ast.FileSpec{
		{
			Name:   "main.go",
			Pkg:    "main",
			Export: true,
			Types: []*ast.TypeSpec{
				{
					Name:   "Config",
					Doc:    "Config doc",
					Export: true,
					Fields: []*ast.FieldSpec{
						{
							Names: []string{"Field1"},
							TypeRef: ast.FieldTypeRef{
								Name: "string",
								Kind: ast.FieldTypeIdent,
							},
							Doc: "Field1 doc",
							Tag: `env:"FIELD1,required,file"`,
						},
					},
				},
				{
					Name:   "Foo",
					Doc:    "Foo doc",
					Export: false,
					Fields: []*ast.FieldSpec{
						{
							Names: []string{"FOne"},
							Doc:   "Foo one field",
							Tag:   `env:"F1"`,
						},
					},
				},
			},
		},
		{
			Name:   "config.go",
			Pkg:    "config",
			Export: false,
			Types: []*ast.TypeSpec{
				{
					Name:   "Bar",
					Doc:    "Bar doc",
					Export: true,
					Fields: []*ast.FieldSpec{
						{
							Names: []string{"BOne"},
							Doc:   "Bar one field",
							Tag:   `env:"B1"`,
						},
					},
				},
			},
		},
	}
	c := NewConverter(types.TargetTypeCaarlos0, opts)
	resolver := resolver.NewTypeResolver()
	scopes := c.ScopesFromFiles(resolver, files)
	expect := []*types.EnvScope{
		{
			Name: "Config",
			Doc:  "Config doc",
			Vars: []*types.EnvDocItem{
				{
					Name: "FIELD1",
					Doc:  "Field1 doc",
					Opts: types.EnvVarOptions{
						Required: true,
						FromFile: true,
					},
				},
			},
		},
	}
	if len(expect) != len(scopes) {
		t.Fatalf("Expected %d scopes, got %d", len(expect), len(scopes))
	}
	for i, scope := range expect {
		checkScope(t, fmt.Sprintf("%d", i), scope, scopes[i])
	}
}

func TestConverterFailedToResolve(t *testing.T) {
	field := &ast.FieldSpec{
		Names: []string{"BarField"},
		TypeRef: ast.FieldTypeRef{
			Pkg:  "config",
			Name: "Bar",
			Kind: ast.FieldTypeIdent,
		},
		Tag: `envPrefix:"BAR_"`,
	}
	c := NewConverter(types.TargetTypeCaarlos0, opts)
	resolver := resolver.NewTypeResolver()
	item := c.DocItemsFromField(resolver, "", field)
	if len(item) != 0 {
		t.Fatalf("Expected 0 items, got %d", len(item))
	}
}

func checkScope(t *testing.T, scope string, expect, actual *types.EnvScope) {
	t.Helper()

	if expect.Name != actual.Name {
		t.Errorf("Expected name %s, got %s", expect.Name, actual.Name)
	}
	if expect.Doc != actual.Doc {
		t.Errorf("Expected doc %s, got %s", expect.Doc, actual.Doc)
	}
	if len(expect.Vars) != len(actual.Vars) {
		t.Fatalf("Expected %d vars, got %d", len(expect.Vars), len(actual.Vars))
	}
	for i, item := range expect.Vars {
		checkDocItem(t, fmt.Sprintf("%s/%d", scope, i), item, actual.Vars[i])
	}
}

func checkDocItem(t *testing.T, scope string, expect, actual *types.EnvDocItem) {
	t.Helper()
	if expect.Name != actual.Name {
		t.Errorf("Expected name %s, got %s", expect.Name, actual.Name)
	}
	if expect.Doc != actual.Doc {
		t.Errorf("Expected doc %s, got %s", expect.Doc, actual.Doc)
	}
	if expect.Opts != actual.Opts {
		t.Errorf("Expected opts %v, got %v", expect.Opts, actual.Opts)
	}
	if len(expect.Children) != len(actual.Children) {
		t.Errorf("Expected %d children, got %d", len(expect.Children), len(actual.Children))
	}
	for i, child := range expect.Children {
		checkDocItem(t, fmt.Sprintf("%s/%d", scope, i), child, actual.Children[i])
	}
}
