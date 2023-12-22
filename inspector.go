package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type inspector struct {
	typeName string // type name to generate documentation for, could be empty
	all      bool   // generate documentation for all types in the file
	execLine int    // line number of the go:generate directive

	lines       []int
	pendingType bool
	items       []*EnvScope
}

func newInspector(typeName string, all bool, execLine int) *inspector {
	return &inspector{typeName: typeName, all: all, execLine: execLine}
}

func (i *inspector) inspectFile(fileName string) (error, []*EnvScope) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err), nil
	}
	// get a lines to position map for the file.
	f := fset.File(file.Pos())
	i.lines = f.Lines()
	items := i.inspect(file)
	return nil, items
}

func (i *inspector) inspect(node ast.Node) []*EnvScope {
	i.items = make([]*EnvScope, 0)
	ast.Walk(i, node)
	return i.items
}

func (i *inspector) getScope(t *ast.TypeSpec) (out *EnvScope) {
	typeName := t.Name.Name
	for _, s := range i.items {
		if s.typeName == typeName {
			out = s
			return
		}
	}

	out = new(EnvScope)
	parseType(t, out)
	i.items = append(i.items, out)
	return
}

func (i *inspector) Visit(n ast.Node) ast.Visitor {
	switch t := n.(type) {
	case *ast.Comment:
		// if type name is not specified we should process the next type
		// declaration after the comment with go:generate
		// which causes this command to be executed.
		if i.typeName != "" || i.all {
			return i
		}
		if !t.Pos().IsValid() {
			return i
		}
		var line int
		for l, pos := range i.lines {
			if token.Pos(pos) > t.Pos() {
				break
			}
			// $GOLINE env var is 1-based.
			line = l + 1
		}
		if line != i.execLine {
			return i
		}

		i.pendingType = true
		return i
	case *ast.TypeSpec:
		var generate bool
		if i.typeName != "" && t.Name != nil && t.Name.Name == i.typeName {
			generate = true
		}
		if i.typeName == "" && i.pendingType {
			generate = true
		}
		if i.all {
			generate = true
		}
		if !generate {
			return i
		}

		if st, ok := t.Type.(*ast.StructType); ok {
			scope := i.getScope(t)
			for _, field := range st.Fields.List {
				var item EnvDocItem
				if !parseField(field, &item) {
					continue
				}
				scope.Vars = append(scope.Vars, item)
			}
		}
		// reset pending type flag event if this type
		// is not processable (e.g. interface type).
		i.pendingType = false
	}
	return i
}

func parseType(t *ast.TypeSpec, out *EnvScope) {
	typeName := t.Name.Name
	out.Doc = strings.TrimSpace(t.Doc.Text())
	out.Name = typeName
	out.typeName = typeName
}

func getTagValues(tag, tagName string) []string {
	tagPrefix := tagName + ":"
	if !strings.Contains(tag, tagPrefix) {
		return nil
	}
	tagValue := strings.Split(tag, tagPrefix)[1]
	leftQ := strings.Index(tagValue, `"`)
	if leftQ == -1 || leftQ == len(tagValue)-1 {
		return nil
	}
	rightQ := strings.Index(tagValue[leftQ+1:], `"`)
	if rightQ == -1 {
		return nil
	}
	tagValue = tagValue[leftQ+1 : leftQ+rightQ+1]
	return strings.Split(tagValue, ",")
}

func parseField(f *ast.Field, out *EnvDocItem) bool {
	if f.Tag == nil {
		return false
	}

	tag := f.Tag.Value
	if !strings.Contains(tag, "env:") {
		return false
	}

	tagValues := getTagValues(tag, "env")
	if len(tagValues) == 0 {
		return false
	}
	out.Name = tagValues[0]
	if f.Doc != nil {
		out.Doc = strings.TrimSpace(f.Doc.Text())
	}
	if out.Doc == "" && f.Comment != nil {
		out.Doc = strings.TrimSpace(f.Comment.Text())
	}

	var opts EnvVarOptions
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

	envDefault := getTagValues(tag, "envDefault")
	if len(envDefault) > 0 {
		opts.Default = strings.Join(envDefault, ",")
	}

	envSeparator := getTagValues(tag, "envSeparator")
	if len(envSeparator) > 0 {
		opts.Separator = envSeparator[0]
	}
	// Check if the field type is a slice or array
	if _, ok := f.Type.(*ast.ArrayType); ok && opts.Separator == "" {
		opts.Separator = ","
	}

	out.Opts = opts

	return true
}
