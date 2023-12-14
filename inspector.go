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
	execLine int

	lines       []int
	pendingType bool
}

func newInspector(typeName string, out inspectorOutput, execLine int) *inspector {
	return &inspector{typeName: typeName, out: out, execLine: execLine}
}

func (i *inspector) inspectFile(fileName string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}
	// get a lines to position map for the file.
	f := fset.File(file.Pos())
	i.lines = f.Lines()
	i.inspect(file)
	return nil
}

func (i *inspector) inspect(node ast.Node) {
	ast.Walk(i, node)
}

func (i *inspector) Visit(n ast.Node) ast.Visitor {
	switch t := n.(type) {
	case *ast.Comment:
		// if type name is not specified we should process the next type
		// declaration after the comment with go:generate
		// which causes this command to be executed.
		if i.typeName != "" {
			return nil
		}
		if !t.Pos().IsValid() {
			return nil
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
			return nil
		}
		if !strings.HasPrefix(t.Text, "//go:generate") {
			return nil
		}

		i.pendingType = true
		return nil
	case *ast.TypeSpec:
		var generate bool
		if i.typeName != "" && t.Name != nil && t.Name.Name == i.typeName {
			generate = true
		}
		if i.typeName == "" && i.pendingType {
			generate = true
		}
		if !generate {
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
		// reset pending type flag event if this type
		// is not processable (e.g. interface type).
		i.pendingType = false
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
