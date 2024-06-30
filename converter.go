package main

import (
	"fmt"
	"os"

	"github.com/g4s8/envdoc/ast"
)

type Converter struct {
	envPrefix     string
	useFieldNames bool
}

func NewConverter(envPrefix string, useFieldNames bool) *Converter {
	return &Converter{
		// envPrefix:     envPrefix,
		useFieldNames: useFieldNames,
	}
}

func (c *Converter) ScopesFromFiles(res *TypeResolver, files []*ast.FileSpec) []*EnvScope {
	var scopes []*EnvScope
	for _, f := range files {
		if !f.Export {
			continue
		}
		for _, t := range f.Types {
			if !t.Export {
				continue
			}
			scopes = append(scopes, c.ScopeFromType(res, t))
		}
	}
	return scopes
}

func (c *Converter) ScopeFromType(res *TypeResolver, t *ast.TypeSpec) *EnvScope {
	scope := &EnvScope{
		Name: t.Name,
		Doc:  t.Doc,
	}
	scope.Vars = c.DocItemsFromFields(res, c.envPrefix, t.Fields)
	return scope
}

func (c *Converter) DocItemsFromFields(res *TypeResolver, prefix string, fields []*ast.FieldSpec) []*EnvDocItem {
	var items []*EnvDocItem
	for _, f := range fields {
		if len(f.Names) == 0 {
			// embedded field
			items = append(items, c.DocItemsFromFields(res, prefix, f.Fields)...)
			continue
		}
		items = append(items, c.DocItemsFromField(res, prefix, f)...)
	}
	return items
}

func (c *Converter) DocItemsFromField(resolver *TypeResolver, prefix string, f *ast.FieldSpec) []*EnvDocItem {
	tag := ParseFieldTag(f.Tag)
	var names []string
	if envName, ok := tag.GetFirst("env"); ok {
		names = []string{envName}
	} else if c.useFieldNames && len(f.Names) > 0 {
		names = make([]string, len(f.Names))
		for i, name := range f.Names {
			names[i] = camelToSnake(name)
		}
	}
	for i, name := range names {
		if name == "" {
			continue
		}
		names[i] = prefix + name
	}

	var opts EnvVarOptions
	if tagValues := tag.GetAll("env"); len(tagValues) > 1 {
		for _, tagValue := range tagValues[1:] {
			switch tagValue {
			case "required":
				opts.Required = true
			case "expand":
				opts.Expand = true
			case "notEmpty":
				opts.Required = true
				opts.NonEmpty = true
			case "file":
				opts.FromFile = true
			}
		}
	}

	if envDefault, ok := tag.GetFirst("envDefault"); ok {
		opts.Default = envDefault
	}

	if envSeparator, ok := tag.GetFirst("envSeparator"); ok {
		opts.Separator = envSeparator
	} else if f.TypeRef.Kind == ast.FieldTypeArray {
		opts.Separator = ","
	}

	var envPrefixed bool
	if envPrefix, ok := tag.GetFirst("envPrefix"); ok {
		prefix = prefix + envPrefix
		envPrefixed = true
	}

	var children []*EnvDocItem
	switch f.TypeRef.Kind {
	case ast.FieldTypeStruct:
		children = c.DocItemsFromFields(resolver, prefix, f.Fields)
	case ast.FieldTypeSelector, ast.FieldTypeIdent, ast.FieldTypeArray, ast.FieldTypePtr:
		if !envPrefixed {
			break
		}
		tpe := resolver.Resolve(&f.TypeRef)
		if tpe == nil {
			fmt.Fprintf(os.Stderr, "Failed to resolve type %q\n", f.TypeRef.String())
			break
		}
		children = c.DocItemsFromFields(resolver, prefix, tpe.Fields)
	}

	res := make([]*EnvDocItem, len(names), len(names)+1)
	for i, name := range names {
		res[i] = &EnvDocItem{
			Name:     name,
			Doc:      f.Doc,
			Opts:     opts,
			Children: children,
		}
	}
	if len(res) == 0 && len(children) > 0 {
		res = append(res, &EnvDocItem{
			Name:     "",
			Doc:      f.Doc,
			Opts:     opts,
			Children: children,
		})
	}
	return res
}
