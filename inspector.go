package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type inspectorOutput interface {
	writeItem(docItem)
}

type inspector struct {
	typeName string
	out      inspectorOutput
}

func newInspector(typeName string, out inspectorOutput) *inspector {
	return &inspector{typeName: typeName, out: out}
}

func (i *inspector) inspectFile(fileName string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}
	i.inspect(file)
	return nil
}

func (i *inspector) inspect(node ast.Node) {
	ast.Walk(i, node)
}

func (i *inspector) Visit(n ast.Node) ast.Visitor {
	switch t := n.(type) {
	case *ast.TypeSpec:
		if t.Name == nil && t.Name.Name != i.typeName {
			return i
		}
		if st, ok := t.Type.(*ast.StructType); ok {
			for _, field := range st.Fields.List {
				if field.Tag != nil {
					var item docItem
					if !parseTag(field.Tag.Value, &item) {
						continue
					}
					item.doc = strings.TrimSpace(field.Doc.Text())
					i.out.writeItem(item)
				}
			}
		}
	}
	return i
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

func parseTag(tag string, out *docItem) bool {
	if !strings.Contains(tag, "env:") {
		return false
	}

	tagValues := getTagValues(tag, "env")
	if len(tagValues) == 0 {
		return false
	}
	out.envName = tagValues[0]
	for _, tagValue := range tagValues[1:] {
		switch tagValue {
		case "required":
			out.flags |= docItemFlagRequired
		case "expand":
			out.flags |= docItemFlagExpand
		case "notEmpty":
			out.flags |= docItemFlagRequired
			out.flags |= docItemFlagNonEmpty
		case "file":
			out.flags |= docItemFlagFromFile
		}
	}

	envDefault := getTagValues(tag, "envDefault")
	if len(envDefault) > 0 {
		out.envDefault = envDefault[0]
	}

	return true
}
