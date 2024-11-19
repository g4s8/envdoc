package linter

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func TestLinter(t *testing.T) {
	for _, tc := range []struct {
		name      string
		file      string
		opts      []Option
		expectOut []string
	}{
		{
			name: "simple",
			file: "testdata/simple.go",
			expectOut: []string{
				"testdata/simple.go:10: field `Undocumented` with `env` tag should have a documentation comment",
				"",
			},
		},
		{
			name: "custom",
			file: "testdata/custom.go",
			opts: []Option{WithEnvName("foo"), WithNoComments()},
			expectOut: []string{
				"testdata/custom.go:10: field `Undocumented` with `foo` tag should have a documentation comment",
				"testdata/custom.go:12: field `NoComments` with `foo` tag should have a documentation comment",
				"",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			log.SetOutput(&out)
			log.SetFlags(0)

			a := NewAnlyzer(false, tc.opts...)

			fset, file := prepareFile(t, tc.file)
			pass := &analysis.Pass{
				Analyzer: a,
				Fset:     fset,
				Files:    []*ast.File{file},
				Report: func(d analysis.Diagnostic) {
					log.Printf("%s:%d: %s", fset.Position(d.Pos).Filename, fset.Position(d.Pos).Line, d.Message)
				},
				ResultOf: make(map[*analysis.Analyzer]interface{}),
			}

			res, err := inspect.Analyzer.Run(pass)
			if err != nil {
				t.Fatalf("could not run inspect analyzer: %v", err)
			}
			pass.ResultOf[inspect.Analyzer] = res

			if _, err := a.Run(pass); err != nil {
				t.Fatalf("could not run linter: %v", err)
			}

			lines := bytes.Split(out.Bytes(), []byte("\n"))
			if len(lines) != len(tc.expectOut) {
				t.Fatalf("unexpected number of lines: got %d, want %d", len(lines), len(tc.expectOut))
			}
			for i, line := range lines {
				if string(line) != tc.expectOut[i] {
					t.Errorf("unexpected output: got %q, want %q", line, tc.expectOut[i])
				}
			}
		})
	}
}

func prepareFile(t *testing.T, name string) (*token.FileSet, *ast.File) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("could not parse file: %v", err)
	}
	return fset, file
}
