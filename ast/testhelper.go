package ast

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
)

type T interface {
	Helper()
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

func loadTestFileSet(t T) (*token.FileSet, *ast.Package, *doc.Package) {
	t.Helper()
	// load go files from ./testdata dir
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "./testdata", nil, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("failed to parse testdata: %s", err)
	}
	pkg, ok := pkgs["testdata"]
	if !ok {
		t.Fatal("package 'testdata' not found")
	}
	docs := doc.New(pkg, "./", doc.PreserveAST|doc.AllDecls)
	return fset, pkg, docs
}

var (
	_ TypeHandler    = (*testTypeHandler)(nil)
	_ CommentHandler = (*testTypeHandler)(nil)
	_ FileHandler    = (*testFileHandler)(nil)
)

type TestTypeHandler interface {
	TypeHandler
	CommentHandler

	Types() []*TypeSpec
}

type testCommentHandler struct {
	comments []*CommentSpec
}

func (h *testCommentHandler) setComment(c *CommentSpec) {
	h.comments = append(h.comments, c)
}

type testSubfieldHandler struct {
	f *FieldSpec
}

func (h *testSubfieldHandler) onField(f *FieldSpec) FieldHandler {
	h.f.Fields = append(h.f.Fields, f)
	return &testSubfieldHandler{f: f}
}

type testFieldHandler struct {
	testCommentHandler
	t *TypeSpec
}

func (h *testFieldHandler) onField(f *FieldSpec) FieldHandler {
	h.t.Fields = append(h.t.Fields, f)
	return &testSubfieldHandler{f: f}
}

type testTypeHandler struct {
	testCommentHandler
	f *FileSpec
}

func (h *testTypeHandler) onType(t *TypeSpec) typeVisitorHandler {
	h.f.Types = append(h.f.Types, t)
	return &testFieldHandler{t: t}
}

func (h *testTypeHandler) Types() []*TypeSpec {
	return h.f.Types
}

type testFileHandler struct {
	testCommentHandler
	files []*FileSpec
}

func (h *testFileHandler) onFile(f *FileSpec) interface {
	TypeHandler
	CommentHandler
} {
	h.files = append(h.files, f)
	return &testTypeHandler{f: f}
}

// func loadTestFileSet(t T) (*token.FileSet, *ast.Package, *doc.Package) {
func testFileVisitor(fset *token.FileSet, pkg *ast.Package, fileName string,
	docs *doc.Package,
) (*testFileHandler, *fileVisitor, *ast.File) {
	fileAst := pkg.Files[fileName]
	fileTkn := fset.File(fileAst.Pos())
	fileSpec := &FileSpec{
		Name: fileTkn.Name(),
		Pkg:  pkg.Name,
	}
	fh := &testFileHandler{}
	th := fh.onFile(fileSpec)
	fv := newFileVisitor(fset, fileAst, docs, th)
	return fh, fv, fileAst
}
