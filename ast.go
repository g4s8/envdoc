package main

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/token"
	"log"
	"strings"
)

//go:generate go run golang.org/x/tools/cmd/stringer@v0.19.0 -type=nodeKind
type nodeKind int

// visitor nodes kinds
const (
	nodeUnknown nodeKind = iota
	nodeType
	nodeRoot
	nodeStruct
	nodeField
)

type visitorNode struct {
	kind        nodeKind
	typeName    string         // type name if node is a type or field type name if node is a field
	packageName string         // package name if node is a type
	currentFile string         // current file name
	names       []string       // it's possible that a field has multiple names
	doc         string         // field or type documentation or comment if doc is empty
	children    []*visitorNode // optional children nodes for structs
	parent      *visitorNode   // parent node
	typeRef     *visitorNode   // type reference if field is a struct
	tag         string         // field tag
	isArray     bool           // true if field is an array
}

type (
	astCommentsHandler func(*ast.Comment) bool
	astTypeDocResolver func(*ast.TypeSpec) string
)

type astVisitor struct {
	commentHandler astCommentsHandler
	logger         *log.Logger
	fileSet        *token.FileSet

	typeDocResolver astTypeDocResolver
	currentNode     *visitorNode
	pendingType     bool   // true if the next type is a target type
	targetName      string // name of the type we are looking for
	depth           int    // current depth in the AST (used for debugging, 1 based)

	err error // error that occurred during AST traversal
}

func newAstVisitor(commentsHandler astCommentsHandler, fileSet *token.FileSet) *astVisitor {
	return &astVisitor{
		commentHandler: commentsHandler,
		fileSet:        fileSet,
		logger:         logger(),
		depth:          1,
	}
}

func (v *astVisitor) push(node *visitorNode, appendChild bool) *astVisitor {
	node.parent = v.currentNode
	if appendChild {
		v.currentNode.children = append(v.currentNode.children, node)
	}
	return &astVisitor{
		commentHandler:  v.commentHandler,
		typeDocResolver: v.typeDocResolver,
		logger:          v.logger,
		pendingType:     v.pendingType,
		currentNode:     node,
		depth:           v.depth + 1,
	}
}

func (v *astVisitor) Walk(n ast.Node) {
	ast.Walk(v, n)
	v.resolveFieldTypes()
}

func (v *astVisitor) Error() error {
	return v.err
}

func (v *astVisitor) setErr(err error) {
	v.logger.Printf("ast(%d): error: %v", v.depth, err)
	v.err = err
}

func (v *astVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if v.err != nil {
		return nil
	}

	v.logger.Printf("ast(%d): visit node (%T)", v.depth, n)

	if v.currentNode == nil {
		v.currentNode = &visitorNode{kind: nodeRoot}
	}

	switch t := n.(type) {
	case *ast.File:
		v.logger.Printf("ast(%d): visit file %q", v.depth, t.Name)
		v.currentNode.packageName = t.Name.Name
		f := v.fileSet.File(t.Pos())
		v.currentNode.currentFile = f.Name()
		typeDocResolver, err := newASTTypeDocResolver(v.fileSet, t)
		if err != nil {
			v.setErr(fmt.Errorf("new ast type doc resolver: %w", err))
			return nil
		}
		v.typeDocResolver = typeDocResolver
		return v
	case *ast.Comment:
		v.logger.Printf("ast(%d): visit comment", v.depth)
		if !v.pendingType {
			v.pendingType = v.commentHandler(t)
		}
		return v
	case *ast.TypeSpec:
		v.logger.Printf("ast(%d): visit type (%T): %q", v.depth, t.Type, t.Name.Name)
		doc := v.typeDocResolver(t)
		name := t.Name.Name
		if v.pendingType {
			v.targetName = name
			v.pendingType = false
			v.logger.Printf("ast(%d): detect target type: %q", v.depth, name)
		}
		typeNode := &visitorNode{
			names:       []string{name},
			typeName:    name,
			packageName: v.currentNode.packageName,
			kind:        nodeType,
			doc:         doc,
		}
		return v.push(typeNode, true)
	case *ast.StructType:
		v.logger.Printf("ast(%d): found struct (`%T`, incomplete: %t, fields: %v)", v.depth, t, t.Incomplete, len(t.Fields.List))
		embedStruct := true
		for i, f := range t.Fields.List {
			if len(f.Names) > 0 {
				embedStruct = false
				break
			}
			v.logger.Printf("ast(%d): debug struct field [%d]%v ", v.depth, i, f.Names)
		}
		if embedStruct {
			v.logger.Printf("ast(%d): struct is embedded", v.depth)
			return v
		}

		switch v.currentNode.kind {
		case nodeType:
			v.currentNode.kind = nodeStruct
			return v
		case nodeField:
			structNode := &visitorNode{
				kind: nodeStruct,
				doc:  v.currentNode.doc,
			}
			v.currentNode.typeRef = structNode
			return v.push(structNode, false)
		default:
			panic(fmt.Sprintf("unexpected node kind: %d", v.currentNode.kind))
		}
	case *ast.Field:
		names := fieldNamesToStr(t)
		v.logger.Printf("ast(%d): visit field ([%d]%v)", v.depth, len(names), names)
		doc := getFieldDoc(t)
		var (
			tag     string
			isArray bool
		)
		if t.Tag != nil {
			tag = t.Tag.Value
		}
		if _, ok := t.Type.(*ast.ArrayType); ok {
			isArray = true
		}
		fieldNode := &visitorNode{
			kind:    nodeField,
			names:   names,
			doc:     doc,
			tag:     tag,
			isArray: isArray,
		}
		if expr, ok := t.Type.(*ast.Ident); ok {
			fieldNode.typeName = expr.Name
		}
		return v.push(fieldNode, true)
	case *ast.FuncDecl:
		return nil
	}
	return v
}

func (v *astVisitor) resolveFieldTypes() {
	unresolved := getAllNodes(v.currentNode, func(n *visitorNode) bool {
		return n.kind == nodeField && n.typeRef == nil
	})
	structs := getAllNodes(v.currentNode, func(n *visitorNode) bool {
		return n.kind == nodeStruct
	})
	structsByName := make(map[string]*visitorNode, len(structs))
	for _, s := range structs {
		structsByName[s.typeName] = s
	}
	for _, f := range unresolved {
		if s, ok := structsByName[f.typeName]; ok {
			f.typeRef = s
			v.logger.Printf("ast: resolve field type %q to struct %q", f.names, s.typeName)
		}
	}
}

func getAllNodes(root *visitorNode, filter func(*visitorNode) bool) []*visitorNode {
	var result []*visitorNode
	if filter(root) {
		result = append(result, root)
	}
	for _, c := range root.children {
		result = append(result, getAllNodes(c, filter)...)
	}
	return result
}

func getFieldDoc(f *ast.Field) string {
	doc := f.Doc.Text()
	if doc == "" {
		doc = f.Comment.Text()
	}
	return strings.TrimSpace(doc)
}

func fieldNamesToStr(f *ast.Field) []string {
	names := make([]string, len(f.Names))
	for i, n := range f.Names {
		names[i] = n.Name
	}
	return names
}

func newASTTypeDocResolver(fileSet *token.FileSet, astFile *ast.File) (func(t *ast.TypeSpec) string, error) {
	docs, err := doc.NewFromFiles(fileSet, []*ast.File{astFile}, "./", doc.PreserveAST|doc.AllDecls)
	if err != nil {
		return nil, fmt.Errorf("extract package docs: %w", err)
	}
	return func(t *ast.TypeSpec) string {
		typeName := t.Name.String()
		docStr := strings.TrimSpace(t.Doc.Text())
		if docStr == "" {
			for _, t := range docs.Types {
				if t.Name == typeName {
					docStr = strings.TrimSpace(t.Doc)
					break
				}
			}
		}
		return docStr
	}, nil
}

var astCommentDummyHandler = func(*ast.Comment) bool {
	return false
}

func newASTCommentTargetLineHandler(goGenLine int, linePositions []int) func(*ast.Comment) bool {
	l := logger()
	return func(c *ast.Comment) bool {
		// if type name is not specified we should process the next type
		// declaration after the comment with go:generate
		// which causes this command to be executed.
		var line int
		for l, pos := range linePositions {
			if token.Pos(pos) > c.Pos() {
				break
			}
			// $GOLINE env var is 1-based.
			line = l + 1
		}
		if line != goGenLine {
			return false
		}

		l.Printf("found go:generate comment at line %d", line)
		return true
	}
}
