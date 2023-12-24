package main

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

type inspector struct {
	typeName      string // type name to generate documentation for, could be empty
	all           bool   // generate documentation for all types in the file
	execLine      int    // line number of the go:generate directive
	useFieldNames bool   // use field names if tag is not specified

	fileSet     *token.FileSet
	lines       []int
	pendingType bool
	items       []*EnvScope
	doc         *doc.Package
	err         error
}

func newInspector(typeName string, all bool, execLine int, useFieldNames bool) *inspector {
	return &inspector{typeName: typeName, all: all, execLine: execLine, useFieldNames: useFieldNames}
}

func (i *inspector) inspectFile(fileName string) ([]*EnvScope, error) {
	i.fileSet = token.NewFileSet()
	file, err := parser.ParseFile(i.fileSet, fileName, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}
	// get a lines to position map for the file.
	f := i.fileSet.File(file.Pos())
	i.lines = f.Lines()
	return i.inspect(file)
}

func (i *inspector) inspect(node ast.Node) ([]*EnvScope, error) {
	i.items = make([]*EnvScope, 0)
	ast.Walk(i, node)
	return i.items, i.err
}

func (i *inspector) getScope(t *ast.TypeSpec) *EnvScope {
	typeName := t.Name.Name
	for _, s := range i.items {
		if s.typeName == typeName {
			return s
		}
	}

	s := i.parseType(t)
	i.items = append(i.items, s)
	return s
}

func (i *inspector) Visit(n ast.Node) ast.Visitor {
	if i.err != nil {
		return nil
	}

	switch t := n.(type) {
	case *ast.File:
		var err error
		i.doc, err = doc.NewFromFiles(i.fileSet, []*ast.File{t}, "./", doc.PreserveAST)
		if err != nil {
			i.err = fmt.Errorf("parse package doc: %w", err)
			return nil
		}
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
				items := i.parseField(field)
				if len(items) == 0 {
					continue
				}
				scope.Vars = append(scope.Vars, items...)
			}
		}
		// reset pending type flag event if this type
		// is not processable (e.g. interface type).
		i.pendingType = false
	}
	return i
}

func (i *inspector) parseType(t *ast.TypeSpec) *EnvScope {
	typeName := t.Name.Name
	docStr := strings.TrimSpace(t.Doc.Text())
	if docStr == "" {
		for _, t := range i.doc.Types {
			if t.Name == typeName {
				docStr = strings.TrimSpace(t.Doc)
				break
			}
		}
	}
	return &EnvScope{
		Name:     typeName,
		Doc:      docStr,
		typeName: typeName,
	}
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

func (i *inspector) parseField(f *ast.Field) (out []EnvDocItem) {
	if f.Tag == nil && !i.useFieldNames {
		return
	}

	var tag string
	if t := f.Tag; t != nil {
		tag = t.Value
	}
	if !strings.Contains(tag, "env:") && !i.useFieldNames {
		return
	}

	tagValues := getTagValues(tag, "env")
	if len(tagValues) > 0 && tagValues[0] != "" {
		var item EnvDocItem
		item.Name = tagValues[0]
		out = []EnvDocItem{item}
	} else if i.useFieldNames {
		out = make([]EnvDocItem, len(f.Names))
		for i, name := range f.Names {
			out[i].Name = camelToSnake(name.Name)
		}
	} else {
		return
	}
	docStr := strings.TrimSpace(f.Doc.Text())
	if docStr == "" {
		docStr = strings.TrimSpace(f.Comment.Text())
	}
	for i := range out {
		out[i].Doc = docStr
	}

	var opts EnvVarOptions
	if len(tagValues) > 1 {
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

	for i := range out {
		out[i].Opts = opts
	}
	return
}
