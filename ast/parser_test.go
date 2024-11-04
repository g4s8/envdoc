package ast

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/g4s8/envdoc/debug"
	"golang.org/x/tools/txtar"
	"gopkg.in/yaml.v2"
)

func TestDataParser(t *testing.T) {
	files, err := filepath.Glob("testdata/parser/*.txtar")
	if err != nil {
		t.Fatalf("failed to list testdata files: %s", err)
	}
	t.Logf("Found %d testdata files", len(files))
	if len(files) == 0 {
		t.Fatal("no testdata files found")
	}

	for _, file := range files {
		file := file

		t.Run(filepath.Base(file), func(t *testing.T) {
			t.Parallel()

			ar, err := txtar.ParseFile(file)
			if err != nil {
				t.Fatalf("failed to parse txtar file: %s", err)
			}

			dir := t.TempDir()
			if err := extractTxtar(ar, dir); err != nil {
				t.Fatalf("failed to extract txtar: %s", err)
			}

			tc := readTestCase(t, dir)
			testParser(t, dir, tc)
		})
	}
}

type parserExpectedTypeRef struct {
	Name    string `yaml:"name"`
	Kind    string `yaml:"kind"`
	Package string `yaml:"pkg"`
}

func (ref *parserExpectedTypeRef) toAST(t *testing.T) FieldTypeRef {
	t.Helper()

	if ref == nil {
		return FieldTypeRef{}
	}

	var kind FieldTypeRefKind
	if !kind.ScanStr(ref.Kind) {
		t.Fatalf("invalid type kind: %s", ref.Kind)
	}
	return FieldTypeRef{
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

func (field *parserExpectedField) toAST(t *testing.T) *FieldSpec {
	t.Helper()

	names := make([]string, len(field.Names))
	copy(names, field.Names)
	fields := make([]*FieldSpec, len(field.Fields))
	for i, f := range field.Fields {
		fields[i] = f.toAST(t)
	}
	return &FieldSpec{
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

func (typ *parserExpectedType) toAST(t *testing.T) *TypeSpec {
	t.Helper()

	fields := make([]*FieldSpec, len(typ.Fields))
	for i, f := range typ.Fields {
		fields[i] = f.toAST(t)
	}
	return &TypeSpec{
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

func (file *parserExpectedFile) toAST(t *testing.T) *FileSpec {
	t.Helper()

	types := make([]*TypeSpec, len(file.Types))
	for i, typ := range file.Types {
		types[i] = typ.toAST(t)
	}
	return &FileSpec{
		Name:   file.Name,
		Pkg:    file.Package,
		Export: file.Exported,
		Types:  types,
	}
}

func parserFilesToAST(t *testing.T, files []*parserExpectedFile) []*FileSpec {
	t.Helper()

	res := make([]*FileSpec, len(files))
	for i, f := range files {
		res[i] = f.toAST(t)
	}
	return res
}

type parserTestCase struct {
	SrcFile  string `yaml:"src_file"`
	FileGlob string `yaml:"file_glob"`
	TypeGlob string `yaml:"type_glob"`
	Debug    bool   `yaml:"debug"`

	Expect []*parserExpectedFile `yaml:"files"`
}

func testParser(t *testing.T, dir string, tc parserTestCase) {
	t.Helper()

	var opts []ParserConfigOption
	if tc.Debug || debug.Config.Enabled {
		debug.SetTestLogger(t)
		t.Log("Debug mode")
		t.Logf("using dir: %s", dir)
		opts = append(opts, WithDebug(true))
	}
	p := NewParser(tc.FileGlob, tc.TypeGlob, opts...)
	files, err := p.Parse(dir)
	if err != nil {
		t.Fatalf("failed to parse files: %s", err)
	}
	astFiles := parserFilesToAST(t, tc.Expect)
	checkFiles(t, "/files", astFiles, files)
}

func checkFiles(t *testing.T, prefix string, expect, res []*FileSpec) {
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

func checkFile(t *testing.T, prefix string, expect, res *FileSpec) {
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

func checkTypes(t *testing.T, prefix string, expect, res []*TypeSpec) {
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

func checkType(t *testing.T, prefix string, expect, res *TypeSpec) {
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

func checkFields(t *testing.T, prefix string, expect, res []*FieldSpec) {
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

func checkField(t *testing.T, prefix string, expect, res *FieldSpec) {
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

func checkTypeRef(t *testing.T, prefix string, expect, res *FieldTypeRef) {
	t.Helper()

	if expect.Name != res.Name {
		t.Errorf("%s: Expected type name %s, got %s", prefix, expect.Name, res.Name)
	}
	if expect.Kind != res.Kind {
		t.Errorf("%s: Expected type kind %s, got %s", prefix, expect.Kind, res.Kind)
	}
}

//---

func extractTxtar(ar *txtar.Archive, dir string) error {
	for _, file := range ar.Files {
		name := filepath.Join(dir, file.Name)
		if err := os.MkdirAll(filepath.Dir(name), 0o777); err != nil {
			return err
		}
		if err := os.WriteFile(name, file.Data, 0o666); err != nil {
			return err
		}
	}
	return nil
}

func readTestCase(t *testing.T, dir string) parserTestCase {
	t.Helper()

	testCaseFile, err := os.Open(path.Join(dir, "testcase.yaml"))
	if err != nil {
		t.Fatalf("failed to open testcase file: %s", err)
	}
	defer testCaseFile.Close()

	var tmp struct {
		TestCase parserTestCase `yaml:"testcase"`
	}
	if err := yaml.NewDecoder(testCaseFile).Decode(&tmp); err != nil {
		t.Fatalf("failed to decode testcase: %s", err)
	}
	return tmp.TestCase
}
