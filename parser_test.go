package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/debug"
	"gopkg.in/yaml.v2"
)

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file: %s", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create destination file: %s", err)
	}
	defer out.Close()

	bufIn := bufio.NewReader(in)
	bufOut := bufio.NewWriter(out)

	if _, err = io.Copy(bufOut, bufIn); err != nil {
		return fmt.Errorf("copy file: %s", err)
	}

	if err = bufOut.Flush(); err != nil {
		return fmt.Errorf("flush destination file: %s", err)
	}
	return nil
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat source dir %q: %w", src, err)
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("create destination dir %q: %w", dst, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read source dir %q: %w", src, err)
	}

	for _, entry := range entries {
		srcPath := path.Join(src, entry.Name())
		dstPath := path.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return fmt.Errorf("copy dir %q: %w", srcPath, err)
			}
			continue
		}
		if err := copyFile(srcPath, dstPath); err != nil {
			return fmt.Errorf("copy file %q: %w", srcPath, err)
		}
	}
	return nil
}

func setupParsserDir(t *testing.T, dir string) string {
	// dir path contains the path to the directory where the test files are located
	// just create a temp dir and copy all files from `dir` to the temp dir recursively
	t.Helper()

	// tmpDir := path.Join(t.TempDir(), "parser")
	tmpDir := path.Join("/tmp/q", "parser")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
	}

	if err := copyDir(path.Join("testdata", dir), tmpDir); err != nil {
		t.Fatalf("failed to copy directory: %s", err)
	}
	return tmpDir
}

func setupParserFiles(t *testing.T, file string) (dir string) {
	t.Helper()

	dir = path.Join(t.TempDir(), "parser")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to create directory: %s", err)
	}
	if err := copyFile(path.Join("testdata", file),
		path.Join(dir, file)); err != nil {
		t.Fatalf("failed to copy file: %s", err)
	}
	return
}

type parserExpectedTypeRef struct {
	Name    string `yaml:"name"`
	Kind    string `yaml:"kind"`
	Package string `yaml:"pkg"`
}

func (ref *parserExpectedTypeRef) toAST(t *testing.T) ast.FieldTypeRef {
	t.Helper()

	if ref == nil {
		return ast.FieldTypeRef{}
	}

	var kind ast.FieldTypeRefKind
	if !kind.ScanStr(ref.Kind) {
		t.Fatalf("invalid type kind: %s", ref.Kind)
	}
	return ast.FieldTypeRef{
		Name: ref.Name,
		Kind: kind,
		Pkg:  ref.Package,
	}
}

type parserExpectedField struct {
	Names   []string               `yaml:"names"`
	Doc     string                 `yaml:"doc"`
	Tag     string                 `yaml:"tag"`
	TypeRef *parserExpectedTypeRef `yaml:"type_ref"`
	Fields  []*parserExpectedField `yaml:"fields"`
}

func (field *parserExpectedField) toAST(t *testing.T) *ast.FieldSpec {
	t.Helper()

	names := make([]string, len(field.Names))
	copy(names, field.Names)
	fields := make([]*ast.FieldSpec, len(field.Fields))
	for i, f := range field.Fields {
		fields[i] = f.toAST(t)
	}
	return &ast.FieldSpec{
		Names:   names,
		Doc:     field.Doc,
		Tag:     field.Tag,
		Fields:  fields,
		TypeRef: field.TypeRef.toAST(t),
	}
}

type parserExpectedType struct {
	Name     string                 `yaml:"name"`
	Exported bool                   `yaml:"export"`
	Doc      string                 `yaml:"doc"`
	Fields   []*parserExpectedField `yaml:"fields"`
}

func (typ *parserExpectedType) toAST(t *testing.T) *ast.TypeSpec {
	t.Helper()

	fields := make([]*ast.FieldSpec, len(typ.Fields))
	for i, f := range typ.Fields {
		fields[i] = f.toAST(t)
	}
	return &ast.TypeSpec{
		Name:   typ.Name,
		Export: typ.Exported,
		Doc:    typ.Doc,
		Fields: fields,
	}
}

type parserExpectedFile struct {
	Name     string                `yaml:"name"`
	Package  string                `yaml:"pkg"`
	Exported bool                  `yaml:"export"`
	Types    []*parserExpectedType `yaml:"types"`
}

func (file *parserExpectedFile) toAST(t *testing.T) *ast.FileSpec {
	t.Helper()

	types := make([]*ast.TypeSpec, len(file.Types))
	for i, typ := range file.Types {
		types[i] = typ.toAST(t)
	}
	return &ast.FileSpec{
		Name:   file.Name,
		Pkg:    file.Package,
		Export: file.Exported,
		Types:  types,
	}
}

func parserFilesToAST(t *testing.T, files []*parserExpectedFile) []*ast.FileSpec {
	t.Helper()

	res := make([]*ast.FileSpec, len(files))
	for i, f := range files {
		res[i] = f.toAST(t)
	}
	return res
}

type parserTestCase struct {
	SrcFile  string `yaml:"src_file"`
	SrcDir   string `yaml:"src_dir"`
	FileGlob string `yaml:"file_glob"`
	TypeGlob string `yaml:"type_glob"`
	Debug    bool   `yaml:"debug"`

	Expect []*parserExpectedFile `yaml:"files"`
}

func loadParserTestCases(t *testing.T) []parserTestCase {
	t.Helper()

	var parserTestCases struct {
		Cases []parserTestCase `yaml:"test_cases"`
	}
	f, err := os.Open("testdata/_cases.yaml")
	if err != nil {
		t.Fatalf("failed to open test cases file: %s", err)
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	dec := yaml.NewDecoder(buf)
	if err := dec.Decode(&parserTestCases); err != nil {
		t.Fatalf("failed to decode test cases: %s", err)
	}
	return parserTestCases.Cases
}

func TestFileParser(t *testing.T) {
	cases := loadParserTestCases(t)
	for i, tc := range cases {
		t.Run(fmt.Sprintf("case(%d)", i), parserTestRunner(tc))
	}
}

func parserTestRunner(tc parserTestCase) func(*testing.T) {
	return func(t *testing.T) {
		var dir string
		if tc.SrcFile != "" {
			dir = setupParserFiles(t, tc.SrcFile)
		} else if tc.SrcDir != "" {
			dir = setupParsserDir(t, tc.SrcDir)
		} else {
			t.Fatal("either src_file or src_dir must be set")
		}
		var opts []parserConfigOption
		if tc.Debug || debug.Config.Enabled {
			debug.SetTestLogger(t)
			t.Log("Debug mode")
			t.Logf("using dir: %s", dir)
			opts = append(opts, withDebug(true))
		}
		p := NewParser(dir, tc.FileGlob, tc.TypeGlob, opts...)
		files, err := p.Parse()
		if err != nil {
			t.Fatalf("failed to parse files: %s", err)
		}
		astFiles := parserFilesToAST(t, tc.Expect)
		checkFiles(t, "/files", astFiles, files)
	}
}

func checkFiles(t *testing.T, prefix string, expect, res []*ast.FileSpec) {
	t.Helper()

	if len(expect) != len(res) {
		t.Errorf("%s: Expected %d files, got %d", prefix, len(expect), len(res))
		for i, file := range expect {
			t.Logf("Expected[%d]: %v", i, file)
		}
		for i, file := range res {
			t.Logf("Got[%d]: %v", i, file)
		}
		return
	}
	for i, file := range expect {
		checkFile(t, fmt.Sprintf("%s/%s", prefix, file.Name), file, res[i])
	}
}

func checkFile(t *testing.T, prefix string, expect, res *ast.FileSpec) {
	t.Helper()

	if !strings.HasSuffix(res.Name, expect.Name) {
		t.Errorf("%s: Expected name %q, got %q", prefix, expect.Name, res.Name)
	}
	if expect.Pkg != res.Pkg {
		t.Errorf("%s: Expected package %q, got %q", prefix, expect.Pkg, res.Pkg)
	}
	if expect.Export != res.Export {
		t.Errorf("%s: Expected export %t, got %t", prefix, expect.Export, res.Export)
	}
	checkTypes(t, prefix+"/types", expect.Types, res.Types)
}

func checkTypes(t *testing.T, prefix string, expect, res []*ast.TypeSpec) {
	t.Helper()

	if len(expect) != len(res) {
		t.Errorf("%s: Expected %d types, got %d", prefix, len(expect), len(res))
		for i, typ := range expect {
			t.Logf("Expected[%d]: %v", i, typ)
		}
		for i, typ := range res {
			t.Logf("Got[%d]: %v", i, typ)
		}
		return
	}
	for i, typ := range expect {
		checkType(t, fmt.Sprintf("%s/%s", prefix, typ.Name), typ, res[i])
	}
}

func checkType(t *testing.T, prefix string, expect, res *ast.TypeSpec) {
	t.Helper()

	if expect.Name != res.Name {
		t.Errorf("%s: Expected name %s, got %s", prefix, expect.Name, res.Name)
	}
	if expect.Doc != res.Doc {
		t.Errorf("%s: Expected doc %s, got %s", prefix, expect.Doc, res.Doc)
	}
	if expect.Export != res.Export {
		t.Errorf("%s: Expected export %t, got %t", prefix, expect.Export, res.Export)
	}
	checkFields(t, prefix+"/fields", expect.Fields, res.Fields)
}

func checkFields(t *testing.T, prefix string, expect, res []*ast.FieldSpec) {
	t.Helper()

	if len(expect) != len(res) {
		t.Errorf("%s: Expected %d fields, got %d", prefix, len(expect), len(res))
		for i, field := range expect {
			t.Logf("Expected[%d]: %s", i, field)
		}
		for i, field := range res {
			t.Logf("Got[%d]: %s", i, field)
		}
		return
	}
	for i, field := range expect {
		str := field.String()
		if str == "" {
			str = fmt.Sprintf("%d", i)
		}
		checkField(t, fmt.Sprintf("%s/%s", prefix, str), field, res[i])
	}
}

func checkField(t *testing.T, prefix string, expect, res *ast.FieldSpec) {
	t.Helper()

	if len(expect.Names) != len(res.Names) {
		t.Errorf("%s: Expected %d names, got %d", prefix, len(expect.Names), len(res.Names))
		for i, name := range expect.Names {
			t.Logf("Expected[%d]: %s", i, name)
		}
		for i, name := range res.Names {
			t.Logf("Got[%d]: %s", i, name)
		}
		return
	}
	for i, name := range expect.Names {
		if name != res.Names[i] {
			t.Errorf("%s: Expected name at %s, got %s", prefix, name, res.Names[i])
		}
	}
	if expect.Doc != res.Doc {
		t.Errorf("%s: Expected doc %s, got %s", prefix, expect.Doc, res.Doc)
	}
	if expect.Tag != res.Tag {
		t.Errorf("%s: Expected tag %q, got %q", prefix, expect.Tag, res.Tag)
	}

	checkTypeRef(t, prefix+"/typeref", &expect.TypeRef, &res.TypeRef)
	checkFields(t, prefix+"/fields", expect.Fields, res.Fields)
}

func checkTypeRef(t *testing.T, prefix string, expect, res *ast.FieldTypeRef) {
	t.Helper()

	if expect.Name != res.Name {
		t.Errorf("%s: Expected type name %s, got %s", prefix, expect.Name, res.Name)
	}
	if expect.Kind != res.Kind {
		t.Errorf("%s: Expected type kind %s, got %s", prefix, expect.Kind, res.Kind)
	}
}
