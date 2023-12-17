package main

import (
	"bufio"
	"bytes"
	"path"
	"testing"
)

func TestGenerator(t *testing.T) {
	testSource := path.Join(t.TempDir(), "test.go")
	if err := copyTestFile(path.Join("testdata", "example_type.go"), testSource); err != nil {
		t.Fatalf("copy test file: %v", err)
	}

	g, err := newGenerator(testSource, 0, withType("Type1"), func(g *generator) error {
		g.format = testFormat(0)
		return nil
	})
	if err != nil {
		t.Fatalf("new generator: %v", err)
	}

	var buf bytes.Buffer
	if err := g.generate(&buf); err != nil {
		t.Fatalf("generate: %v", err)
	}

	expect := []string{
		"#BEFORE",
		"FOO|Foo stub",
		"#AFTER",
	}
	scanner := bufio.NewScanner(&buf)
	var count int
	for scanner.Scan() {
		if count >= len(expect) {
			t.Fatalf("unexpected line: %q", scanner.Text())
		}
		if scanner.Text() != expect[count] {
			t.Fatalf("unexpected line: %q", scanner.Text())
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner: %v", err)
	}
	if count != len(expect) {
		t.Fatalf("unexpected line count: %d", count)
	}
}

func TestGeneratorOpts(t *testing.T) {
	t.Run("withType", func(t *testing.T) {
		g, err := newGenerator("", 0, withType("Type1"), func(g *generator) error {
			g.format = testFormat(0)
			return nil
		})
		if err != nil {
			t.Fatalf("new generator: %v", err)
		}
		if targetType := g.targetType; targetType != "Type1" {
			t.Fatalf("unexpected type name: %q", targetType)
		}
	})
	t.Run("withFormat", func(t *testing.T) {
		t.Run("invalid", func(t *testing.T) {
			_, err := newGenerator("", 0, withFormat("invalid"))
			if err == nil {
				t.Fatalf("expected error")
			}
		})
		t.Run("markdown", func(t *testing.T) {
			g, err := newGenerator("", 0, withFormat("markdown"))
			if err != nil {
				t.Fatalf("new generator: %v", err)
			}
			if format := g.format; format != fmtMD {
				t.Fatalf("unexpected format: %v", format)
			}
		})
		t.Run("plaintext", func(t *testing.T) {
			g, err := newGenerator("", 0, withFormat("plaintext"))
			if err != nil {
				t.Fatalf("new generator: %v", err)
			}
			if format := g.format; format != fmtPlain {
				t.Fatalf("unexpected format: %v", format)
			}
		})
		t.Run("html", func(t *testing.T) {
			g, err := newGenerator("", 0, withFormat("html"))
			if err != nil {
				t.Fatalf("new generator: %v", err)
			}
			if format := g.format; format != fmtHTML {
				t.Fatalf("unexpected format: %v", format)
			}
		})
	})
}
