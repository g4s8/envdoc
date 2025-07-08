package ast

import (
	"go/ast"
	"go/doc"
	"go/token"
	"strings"

	"github.com/g4s8/envdoc/debug"
)

func getFieldTypeRef(f ast.Expr, ref *FieldTypeRef) bool {
	switch t := f.(type) {
	case *ast.Ident:
		ref.Name = t.Name
		ref.Kind = FieldTypeIdent
	case *ast.SelectorExpr:
		getFieldTypeRef(t.X, ref)
		ref.Pkg = ref.Name
		ref.Name = t.Sel.Name
		ref.Kind = FieldTypeSelector
	case *ast.StarExpr:
		getFieldTypeRef(t.X, ref)
		ref.Kind = FieldTypePtr
	case *ast.ArrayType:
		getFieldTypeRef(t.Elt, ref)
		ref.Kind = FieldTypeArray
	case *ast.MapType:
		getFieldTypeRef(t.Value, ref)
		ref.Kind = FieldTypeMap
	case *ast.StructType:
		ref.Kind = FieldTypeStruct
	default:
		return false
	}
	return true
}

func extractFieldNames(f *ast.Field) []string {
	names := make([]string, len(f.Names))
	for i, n := range f.Names {
		names[i] = n.Name
	}
	return names
}

func extractFieldDoc(f *ast.Field) (doc string, ok bool) {
	doc = f.Doc.Text()
	if doc == "" {
		doc = f.Comment.Text()
	}
	doc = strings.TrimSpace(doc)
	return doc, doc != ""
}

func findCommentLine(c *ast.Comment, fset *token.FileSet, file *ast.File) int {
	lines := fset.File(file.Pos()).Lines()
	for l, pos := range lines {
		if token.Pos(pos) > c.Pos() {
			return l
		}
	}
	return 0
}

func getFieldSpec(n *ast.Field, pkg string) *FieldSpec {
	names := extractFieldNames(n)
	allPrivate := true
	for _, name := range names {
		if strings.ToLower(name[:1]) != name[:1] {
			allPrivate = false
			break
		}
	}
	if len(names) > 0 && allPrivate {
		// skip private fields
		return nil
	}

	var fs FieldSpec
	fs.Names = names
	if !getFieldTypeRef(n.Type, &fs.TypeRef) {
		// unsupported field type
		return nil
	}
	if fs.TypeRef.Pkg == "" {
		fs.TypeRef.Pkg = pkg
	}
	if doc, ok := extractFieldDoc(n); ok {
		fs.Doc = doc
	}
	if tag := n.Tag; tag != nil {
		fs.Tag = strings.Trim(tag.Value, "`")
	}

	return &fs
}

//nolint:cyclop
func debugNode(src string, n ast.Node) {
	if !debug.Config.Enabled {
		return
	}
	if n == nil {
		return
	}

	switch t := n.(type) {
	case *ast.File:
		debug.Logf("# AST(%s): File pkg=%q\n", src, t.Name.Name)
	case *ast.ImportSpec:
		debug.Logf("# AST(%s): Import %s %s\n", src, t.Name, t.Path.Value)
	//nolint:staticcheck
	case *ast.Package:
		debug.Logf("# AST(%s): Package %s\n", src, t.Name)
	case *ast.TypeSpec:
		debug.Logf("# AST(%s): Type %s\n", src, t.Name.Name)
	case *ast.Field:
		names := extractFieldNames(t)
		debug.Logf("# AST(%s): Field %s\n", src, strings.Join(names, ", "))
	case *ast.Comment:
		debug.Logf("# AST(%s): Comment %s\n", src, t.Text)
	case *ast.StructType:
		debug.Logf("# AST(%s): Struct\n", src)
	case *ast.GenDecl, *ast.Ident, *ast.FuncDecl:
		// ignore
	default:
		debug.Logf("# AST(%s): %T\n", src, t)
	}
}

func resolveTypeDocs(docs *doc.Package, t *ast.TypeSpec) string {
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
}
