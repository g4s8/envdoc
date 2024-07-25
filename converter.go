package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/debug"
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
			debug.Logf("# CONV: skip file %q\n", f.Name)
			continue
		}
		for _, t := range f.Types {
			if !t.Export {
				debug.Logf("# CONV: skip type %q\n", t.Name)
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
	debug.Logf("# CONV: found scope %q\n", scope.Name)
	return scope
}

func (c *Converter) DocItemsFromFields(res *TypeResolver, prefix string, fields []*ast.FieldSpec) []*EnvDocItem {
	var items []*EnvDocItem
	for _, f := range fields {
		debug.Logf("\t# CONV: field [%s]\n", strings.Join(f.Names, ","))
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
		// if name == "" {
		// 	continue
		// }
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
		debug.Logf("\t# CONV: struct %q (%d childrens)\n", f.TypeRef.String(), len(children))
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
		debug.Logf("\t# CONV: selector %q (%d childrens)\n", f.TypeRef.String(), len(children))
	}

	res := make([]*EnvDocItem, len(names), len(names)+1)
	for i, name := range names {
		res[i] = &EnvDocItem{
			Name:     name,
			Doc:      f.Doc,
			Opts:     opts,
			Children: children,
		}
		debug.Logf("\t# CONV: docItem %q (%d childrens)\n", name, len(children))
	}

	if len(names) == 0 && len(children) > 0 {
		return children
	}
	return res
}
