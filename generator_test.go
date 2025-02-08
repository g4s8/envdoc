package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/render"
	"github.com/g4s8/envdoc/types"
	"golang.org/x/tools/txtar"
)

func TestGenerator(t *testing.T) {
	files, err := filepath.Glob("testdata/*.txtar")
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
			spec := parseTestSpec(t, string(ar.Comment))
			t.Logf("Test case: %s", spec.Comment)

			dir := extractTxtar(t, ar)

			p := ast.NewParser("*", spec.TypeName)
			conv := NewConverter(types.TargetTypeCaarlos0, ConverterOpts{
				EnvPrefix:     spec.EnvPrefix,
				TagName:       "env",
				TagDefault:    "envDefault",
				UseFieldNames: spec.FieldNames,
			})
			rend := render.NewRenderer(types.OutFormatTxt, false)
			gen := NewGenerator(p, conv, rend)
			var out bytes.Buffer
			runGenerator(t, gen, spec, dir, &out)

			expectFile, err := os.Open(path.Join(dir, "expect.txt"))
			if err != nil {
				t.Fatalf("failed to open expect.txt: %s", err)
			}
			defer expectFile.Close()
			expect, err := io.ReadAll(expectFile)
			if err != nil {
				t.Fatalf("failed to read expect.txt: %s", err)
			}
			if !bytes.Equal(out.Bytes(), expect) {
				t.Logf("Expected:\n%s", expect)
				t.Logf("Got:\n%s", out.String())
				t.Fatalf("Output mismatch")
			}
		})
	}
}

func extractTxtar(t *testing.T, ar *txtar.Archive) string {
	t.Helper()

	dir := t.TempDir()
	for _, file := range ar.Files {
		name := filepath.Join(dir, file.Name)
		if err := os.MkdirAll(filepath.Dir(name), 0o777); err != nil {
			t.Fatalf("failed to create dir: %s", err)
		}
		if err := os.WriteFile(name, file.Data, 0o666); err != nil {
			t.Fatalf("failed to write file: %s", err)
		}
	}
	return dir
}

func runGenerator(t *testing.T, gen interface{ Generate(string, io.Writer) error }, spec GenTestSpec, dir string, out *bytes.Buffer) {
	t.Helper()

	if err := gen.Generate(dir, out); err != nil {
		if spec.Success {
			t.Fatalf("failed to generate: %s", err)
		}
		t.Logf("Expected error: %s", err)
		return
	}
	if !spec.Success {
		t.Fatalf("expected error, but got success")
	}
}

type GenTestSpec struct {
	Success    bool
	TypeName   string
	EnvPrefix  string
	FieldNames bool
	Comment    string
}

func parseTestSpec(t *testing.T, data string) GenTestSpec {
	t.Helper()

	var res GenTestSpec
	// comment is a multiline test spec.
	// the first line starts with either `Success:` or `Error`
	// with following description of test case.
	// If the first line starts with `Error`, the test is expected to fail.
	// Next lines may contain:
	// - TypeName: type name to process
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Success:") {
			res.Success = true
			continue
		}
		if strings.HasPrefix(line, "Success:") || strings.HasPrefix(line, "Error:") {
			s := strings.SplitN(line, ":", 2)
			res.Comment = strings.TrimSpace(s[1])
			continue
		}

		if strings.HasPrefix(line, "TypeName:") {
			res.TypeName = strings.TrimSpace(strings.TrimPrefix(line, "TypeName:"))
			continue
		}
		if strings.HasPrefix(line, "EnvPrefix:") {
			res.EnvPrefix = strings.TrimSpace(strings.TrimPrefix(line, "EnvPrefix:"))
			continue
		}
		if strings.HasPrefix(line, "FieldNames:") {
			res.FieldNames = strings.TrimSpace(strings.TrimPrefix(line, "FieldNames:")) == "true"
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("failed to read comment: %s", err)
	}
	if res.TypeName == "" {
		res.TypeName = "*"
	}
	return res
}
