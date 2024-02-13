package main

import (
	"bytes"
	"errors"
	"path"
	"testing"
)

func TestOptions(t *testing.T) {
	t.Run("WithFormat", func(t *testing.T) {
		for _, c := range []struct {
			name      string
			expect    any
			expectErr bool
		}{
			{name: "", expect: tmplMarkdown},
			{name: "markdown", expect: tmplMarkdown},
			{name: "html", expect: tmplHTML},
			{name: "plaintext", expect: tmplPlaintext},
			{name: "unknown", expectErr: true},
		} {
			t.Run(c.name, func(t *testing.T) {
				g, err := newGenerator("stub", 1, withFormat(c.name))
				if err != nil && !c.expectErr {
					t.Fatal("new generator error", err)
				}
				if err == nil && c.expectErr {
					t.Fatal("expected error, got nil")
				}
				if !c.expectErr && g.tmpl != c.expect {
					t.Errorf("expected %v, got %v", c.expect, g.tmpl)
				}
			})
		}
	})
	t.Run("WithAll", func(t *testing.T) {
		g, err := newGenerator("stub", 1, withAll())
		if err != nil {
			t.Fatal("new generator error", err)
		}
		if g.all != true {
			t.Fatal("expected all to be true")
		}
	})
	t.Run("WithPrefix", func(t *testing.T) {
		g, err := newGenerator("stub", 1, withPrefix("prefix"))
		if err != nil {
			t.Fatal("new generator error", err)
		}
		if g.prefix != "prefix" {
			t.Fatal("expected prefix to be 'prefix'")
		}
	})
	t.Run("WithNoStyles", func(t *testing.T) {
		g, err := newGenerator("stub", 1, withNoStyles())
		if err != nil {
			t.Fatal("new generator error", err)
		}
		if g.noStyles != true {
			t.Fatal("expected noStyles to be true")
		}
	})
	t.Run("WithFieldNames", func(t *testing.T) {
		g, err := newGenerator("stub", 1, withFieldNames())
		if err != nil {
			t.Fatal("new generator error", err)
		}
		if g.fieldNames != true {
			t.Fatal("expected fieldNames to be true")
		}
	})
	t.Run("WithType", func(t *testing.T) {
		const typeName = "Foo"
		g, err := newGenerator("stub", 1, withType(typeName))
		if err != nil {
			t.Fatal("new generator error", err)
		}
		if g.targetType != typeName {
			t.Fatalf("expected targetType to be %q, got %q", typeName, g.targetType)
		}
	})
	t.Run("empty", func(t *testing.T) {
		g, err := newGenerator("stub", 1)
		if err != nil {
			t.Fatal("create new generator", err)
		}
		if tmpl := g.tmpl; tmpl != tmplMarkdown {
			t.Fatal("invalid default template", tmpl)
		}
	})
}

type brokenWriter struct{}

var errBroken = errors.New("broken")

func (w brokenWriter) Write(p []byte) (n int, err error) {
	return 0, errBroken
}

func TestGenerator(t *testing.T) {
	t.Run("broken-input", func(t *testing.T) {
		g, err := newGenerator("stub", 1, withFormat("markdown"))
		if err != nil {
			t.Fatal("new generator error", err)
		}
		var out bytes.Buffer
		err = g.generate(&out)
		if err == nil {
			t.Error("expected error, got nil")
		}
		t.Logf("got expected error: %v", err)
	})
	t.Run("broken-out", func(t *testing.T) {
		src := path.Join(t.TempDir(), "example.go")
		if err := copyTestFile("testdata/tags.go", src); err != nil {
			t.Fatal("copy test file error", err)
		}
		g, err := newGenerator(src, 1, withFormat(""))
		if err != nil {
			t.Fatal("new generator error", err)
		}
		err = g.generate(brokenWriter{})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
	t.Run("success", func(t *testing.T) {
		src := path.Join(t.TempDir(), "example.go")
		if err := copyTestFile("testdata/tags.go", src); err != nil {
			t.Fatal("copy test file error", err)
		}
		g, err := newGenerator(src, 1, withFormat(""))
		if err != nil {
			t.Fatal("new generator error", err)
		}
		var out bytes.Buffer
		if err := g.generate(&out); err != nil {
			t.Fatal("generate error", err)
		}
		if out.Len() == 0 {
			t.Error("expected output, got empty")
		}
	})
}
