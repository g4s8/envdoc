package main

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

type envFieldKind int

const (
	envFieldKindPlain  envFieldKind = iota
	envFieldKindStruct              // struct reference
)

type envField struct {
	name      string
	kind      envFieldKind
	doc       string
	opts      EnvVarOptions
	typeRef   string
	envPrefix string
}

type envStruct struct {
	name   string
	doc    string
	fields []envField
}

type inspector struct {
	typeName      string // type name to generate documentation for, could be empty
	all           bool   // generate documentation for all types in the file
	execLine      int    // line number of the go:generate directive
	useFieldNames bool   // use field names if tag is not specified

	fileSet     *token.FileSet
	lines       []int
	pendingType bool
	items       []*envStruct
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
	i.items = make([]*envStruct, 0)
	ast.Walk(i, node)
	if i.err != nil {
		return nil, i.err
	}
	scopes, err := i.buildScopes()
	if err != nil {
		return nil, fmt.Errorf("build scopes: %w", err)
	}
	return scopes, nil
}

func (i *inspector) getStruct(t *ast.TypeSpec) *envStruct {
	typeName := t.Name.Name
	for _, s := range i.items {
		if s.name == typeName {
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
		if i.typeName == "" && i.pendingType {
			i.typeName = t.Name.Name
		}

		if st, ok := t.Type.(*ast.StructType); ok {
			str := i.getStruct(t)
			debug("parsing struct %s", str.name)
			for _, field := range st.Fields.List {
				items := i.parseField(field)
				for i, item := range items {
					if item.kind == envFieldKindPlain {
						debug("parsed field[%d] %s", i, item.name)
					} else {
						debug("parsed field[%d] %s (struct ref: %s, prefix: %s)", i, item.name, item.typeRef, item.envPrefix)
					}
				}
				if len(items) == 0 {
					continue
				}
				str.fields = append(str.fields, items...)
			}
		}
		// reset pending type flag event if this type
		// is not processable (e.g. interface type).
		i.pendingType = false
	}
	return i
}

func (i *inspector) parseType(t *ast.TypeSpec) *envStruct {
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
	return &envStruct{
		name: typeName,
		doc:  docStr,
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

func (i *inspector) parseField(f *ast.Field) (out []envField) {
	if f.Tag == nil && !i.useFieldNames {
		return
	}

	var tag string
	if t := f.Tag; t != nil {
		tag = t.Value
	}

	envPrefix := getTagValues(tag, "envPrefix")
	if len(envPrefix) > 0 && envPrefix[0] != "" {
		var item envField
		item.envPrefix = envPrefix[0]
		item.kind = envFieldKindStruct
		fieldType := f.Type.(*ast.Ident)
		item.typeRef = fieldType.Name
		out = []envField{item}
		return
	}

	if !strings.Contains(tag, "env:") && !i.useFieldNames {
		return
	}

	tagValues := getTagValues(tag, "env")
	if len(tagValues) > 0 && tagValues[0] != "" {
		var item envField
		item.name = tagValues[0]
		item.kind = envFieldKindPlain
		out = []envField{item}
	} else if i.useFieldNames {
		out = make([]envField, len(f.Names))
		for i, name := range f.Names {
			out[i].name = camelToSnake(name.Name)
			out[i].kind = envFieldKindPlain
		}
	} else {
		return
	}
	docStr := strings.TrimSpace(f.Doc.Text())
	if docStr == "" {
		docStr = strings.TrimSpace(f.Comment.Text())
	}
	for i := range out {
		out[i].doc = docStr
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
		out[i].opts = opts
	}
	return
}

func (i *inspector) buildScopes() ([]*EnvScope, error) {
	scopes := make([]*EnvScope, 0, len(i.items))
	for _, s := range i.items {
		if !i.all && s.name != i.typeName {
			debug("skip %q", s.name)
			continue
		}

		debug("process %q", s.name)
		scope := &EnvScope{
			Name: s.name,
			Doc:  s.doc,
		}
		for _, f := range s.fields {
			switch f.kind {
			case envFieldKindPlain:
				v := EnvDocItem{
					Name: f.name,
					Doc:  f.doc,
					Opts: f.opts,
				}
				debug("[p] add docItem: %s <- %s", scope.Name, v.Name)
				scope.Vars = append(scope.Vars, v)
			case envFieldKindStruct:
				envPrefix := f.envPrefix
				var base *envStruct
				for _, s := range i.items {
					if s.name == f.typeRef {
						base = s
						break
					}
				}
				if base == nil {
					return nil, fmt.Errorf("struct %q not found", f.typeRef)
				}
				for _, f := range base.fields {
					name := fmt.Sprintf("%s%s", envPrefix, f.name)
					v := EnvDocItem{
						Name: name,
						Doc:  f.doc,
						Opts: f.opts,
					}
					debug("[s] add docItem: %s <- %s (prefix: %s)", scope.Name, v.Name, envPrefix)
					scope.Vars = append(scope.Vars, v)
				}
			default:
				panic("unknown field kind")
			}
		}
		scopes = append(scopes, scope)
	}
	return scopes, nil
}

const debugLogs = false

func debug(f string, args ...any) {
	if !debugLogs {
		return
	}
	fmt.Printf("DEBUG: "+f+"\n", args...)
}
