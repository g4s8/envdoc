package main

import (
	"fmt"
	"testing"

	"github.com/g4s8/envdoc/ast"
)

func TestConvertDocItems(t *testing.T) {
	c := NewConverter("", true)
	fieldValues := []*ast.FieldSpec{
		{
			Names: []string{"Field1"},
			TypeRef: ast.FieldTypeRef{
				Name: "string",
				Kind: ast.FieldTypeIdent,
			},
			Tag: `env:"FIELD1,required"`,
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
					Tag: `env:"FIELD4,notEmpty"`,
				},
			},
			Tag: `envPrefix:"PREFIX_"`,
			TypeRef: ast.FieldTypeRef{
				Kind: ast.FieldTypeStruct,
			},
		},
	}
	resolver := NewTypeResolver()
	res := c.DocItemsFromFields(resolver, "", fieldValues)
	if len(res) != 4 {
		t.Errorf("Expected 4 items, got %d", len(res))
	}
	expect := []*EnvDocItem{
		{
			Name: "FIELD1",
			Doc:  "Field1 doc",
			Opts: EnvVarOptions{
				Required: true,
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
			Name: "FIELD4",
			Doc:  "Field4 doc",
			Opts: EnvVarOptions{
				Required: true,
				NonEmpty: true,
			},
		},
	}
	for i, item := range expect {
		checkDocItem(t, fmt.Sprintf("%d", i), item, res[i])
	}
}

func checkDocItem(t *testing.T, scope string, expect, actual *EnvDocItem) {
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
