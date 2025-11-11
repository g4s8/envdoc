package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/debug"
	"github.com/g4s8/envdoc/types"
)

type Resolver interface {
	Resolve(f *ast.FileSpec, ref *ast.FieldTypeRef) *ast.TypeSpec
}

type ConverterOpts struct {
	EnvPrefix       string
	TagName         string
	TagDefault      string
	RequiredIfNoDef bool
	UseFieldNames   bool
	CustomTemplate  string
}

type Converter struct {
	target types.TargetType
	opts   ConverterOpts
}

func NewConverter(target types.TargetType, opts ConverterOpts) *Converter {
	return &Converter{
		target: target,
		opts:   opts,
	}
}

func (c *Converter) ScopesFromFiles(res Resolver, files []*ast.FileSpec) []*types.EnvScope {
	var scopes []*types.EnvScope
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
			scopes = append(scopes, c.ScopeFromType(res, f, t))
		}
	}
	return scopes
}

func (c *Converter) ScopeFromType(res Resolver, file *ast.FileSpec, t *ast.TypeSpec) *types.EnvScope {
	scope := &types.EnvScope{
		Name: t.Name,
		Doc:  t.Doc,
	}
	scope.Vars = c.DocItemsFromFields(res, file, c.opts.EnvPrefix, t.Fields)
	debug.Logf("# CONV: found scope %q\n", scope.Name)
	return scope
}

func (c *Converter) DocItemsFromFields(res Resolver, file *ast.FileSpec, prefix string, fields []*ast.FieldSpec) []*types.EnvDocItem {
	var items []*types.EnvDocItem
	for _, f := range fields {
		debug.Logf("\t# CONV: field [%s] type=%s flen=%d\n",
			strings.Join(f.Names, ","), f.TypeRef, len(f.Fields))
		if len(f.Names) == 0 {
			// embedded field
			if len(f.Fields) == 0 {
				// resolve embedded types
				tpe := res.Resolve(file, &f.TypeRef)
				if tpe != nil {
					f.Fields = tpe.Fields
				}
			}
			items = append(items, c.DocItemsFromFields(res, file, prefix, f.Fields)...)
			continue
		}
		items = append(items, c.DocItemsFromField(res, file, prefix, f)...)
	}
	return items
}

func (c *Converter) DocItemsFromField(resolver Resolver, file *ast.FileSpec, prefix string, f *ast.FieldSpec) []*types.EnvDocItem {
	dec := NewFieldDecoder(c.target, FieldDecoderOpts{
		EnvPrefix:       prefix,
		TagName:         c.opts.TagName,
		TagDefault:      c.opts.TagDefault,
		RequiredIfNoDef: c.opts.RequiredIfNoDef,
		UseFieldNames:   c.opts.UseFieldNames,
	})
	info, newPrefix := dec.Decode(f)
	if newPrefix != "" {
		prefix = newPrefix
	}

	var children []*types.EnvDocItem
	switch f.TypeRef.Kind {
	case ast.FieldTypeStruct:
		children = c.DocItemsFromFields(resolver, file, prefix, f.Fields)
		debug.Logf("\t# CONV: struct %q (%d childrens)\n", f.TypeRef.String(), len(children))
	case ast.FieldTypeSelector, ast.FieldTypeIdent, ast.FieldTypeArray, ast.FieldTypePtr:
		if f.TypeRef.IsBuiltIn() {
			break
		}
		tpe := resolver.Resolve(file, &f.TypeRef)
		debug.Logf("\t# CONV: resolve %q -> %v\n", f.TypeRef.String(), tpe)
		if tpe == nil {
			if newPrefix != "" {
				// Target type is env-prefixed, it means it's a reference
				// to another struct type. We can't process it here, because
				// we can't resolve the target type and its fields.
				fmt.Fprintf(os.Stderr, "WARNING: failed to resolve type %q\n", f.TypeRef.String())
			}
			break
		}
		children = c.DocItemsFromFields(resolver, file, prefix, tpe.Fields)
		debug.Logf("\t# CONV: selector %q (%d childrens)\n", f.TypeRef.String(), len(children))
	}

	res := make([]*types.EnvDocItem, len(info.Names), len(info.Names)+1)
	opts := types.EnvVarOptions{
		Required:  info.Required,
		Expand:    info.Expand,
		NonEmpty:  info.NonEmpty,
		FromFile:  info.FromFile,
		Default:   info.Default,
		Separator: info.Separator,
	}
	for i, name := range info.Names {
		res[i] = &types.EnvDocItem{
			Name:     name,
			Doc:      f.Doc,
			Opts:     opts,
			Children: children,
		}
		debug.Logf("\t# CONV: docItem %q (%d childrens)\n", name, len(children))
	}

	if len(info.Names) == 0 && len(children) > 0 {
		return children
	}
	return res
}

func (c *Converter) CustomTemplate() string {
	return c.opts.CustomTemplate
}
