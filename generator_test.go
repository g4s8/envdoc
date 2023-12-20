package main

import (
	"bytes"
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
	t.Run("empty", func(t *testing.T) {
		_, err := newGenerator("stub", 1)
		if err == nil {
			t.Error("expected error, got nil")
		}
		t.Logf("got expected error: %v", err)
	})
}

func TestGenerator(t *testing.T) {
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
}
