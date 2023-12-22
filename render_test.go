package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

var testRenderItems = []renderItem{
	{
		EnvName: "TEST_ENV",
		Doc:     "This is a test environment variable.",
	},
	{
		EnvName:      "TEST_ENV2",
		Doc:          "This is another test environment variable.",
		EnvDefault:   "default value",
		EnvSeparator: ",",
	},
	{
		EnvName:  "TEST_ENV3",
		Doc:      "This is a third test environment variable.",
		Required: true,
		Expand:   true,
		NonEmpty: true,
		FromFile: true,
	},
}

var testRenderSections = []renderSection{
	{
		Name: "First",
		Items: []renderItem{
			{
				EnvName: "ONE",
				Doc:     "First one",
			}, {
				EnvName: "TWO",
				Doc:     "First two",
			},
		},
	},
	{
		Name: "Second",
		Items: []renderItem{
			{
				EnvName: "THREE",
				Doc:     "Second three",
			},
		},
	},
}

func TestRender(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		rc := renderContext{Title: "Simple", Sections: []renderSection{{Items: testRenderItems}}}
		t.Run("markdown", testRenderer(tmplMarkdown, rc,
			"# Simple",
			"- `TEST_ENV` - This is a test environment variable.",
			"- `TEST_ENV2` (comma-separated, default: `default value`) - This is another test environment variable.",
			"- `TEST_ENV3` (**required**, expand, non-empty, from-file) - This is a third test environment variable."))
		t.Run("plaintext", testRenderer(tmplPlaintext, rc,
			"Simple",
			" * `TEST_ENV` - This is a test environment variable.",
			" * `TEST_ENV2` (comma-separated, default: `default value`) - This is another test environment variable.",
			" * `TEST_ENV3` (required, expand, non-empty, from-file) - This is a third test environment variable."))
		t.Run("html", testRenderer(tmplHTML, rc,
			`<!DOCTYPE html>`,
			`<html lang="en">`,
			`<head>`,
			`<meta charset="utf-8">`,
			`<title>Simple</title>`,
			`</head>`,
			`<section>`,
			`<article>`,
			`<h1>Simple</h1>`,
			`<ul>`,
			`<li><code>TEST_ENV</code> - This is a test environment variable.</li>`,
			`<li><code>TEST_ENV2</code> (comma-separated, default: <code>default value</code>) - This is another test environment variable.</li>`,
			`<li><code>TEST_ENV3</code> (<strong>required</strong>, expand, non-empty, from-file) - This is a third test environment variable.</li>`,
			`</ul>`,
			`</article>`,
			`</section>`,
			`</body>`,
			`</html>`))
	})
	t.Run("sections", func(t *testing.T) {
		rc := renderContext{Title: "Sections", Sections: testRenderSections}
		t.Run("markdown", testRenderer(tmplMarkdown, rc,
			"# Sections",
			"## First",
			" - `ONE` - First one",
			" - `TWO` - First two",
			"## Second",
			" - `THREE` - Second three"))
		t.Run("plaintext", testRenderer(tmplPlaintext, rc,
			"Sections",
			"## First",
			" * `ONE` - First one",
			" * `TWO` - First two",
			"## Second",
			" * `THREE` - Second three"))
		t.Run("html", testRenderer(tmplHTML, rc,
			`<!DOCTYPE html>`,
			`<html lang="en">`,
			`<head>`,
			`<meta charset="utf-8">`,
			`<title>Sections</title>`,
			`</head>`,
			`<section>`,
			`<article>`,
			`<h1>Sections</h1>`,
			`<h2>First</h2>`,
			`<ul>`,
			`<li><code>ONE</code> - First one</li>`,
			`<li><code>TWO</code> - First two</li>`,
			`</ul>`,
			`<h2>Second</h2>`,
			`<li><code>THREE</code> - Second three</li>`,
			`</ul>`,
			`</article>`,
			`</section>`,
			`</body>`,
			`</html>`))
	})
}

func TestNewRenderContext(t *testing.T) {
	src := []*EnvScope{
		{
			Name: "First",
			Vars: []EnvDocItem{
				{
					Name: "ONE",
					Doc:  "First one",
				},
			},
		},
	}
	rc := newRenderContext(src, "PREFIX_")
	const title = "Environment Variables"
	if rc.Title != title {
		t.Errorf("expected title %q, got %q", title, rc.Title)
	}
	if len(rc.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(rc.Sections))
	}
	section := rc.Sections[0]
	if section.Name != "First" {
		t.Errorf("expected section name %q, got %q", "First", section.Name)
	}
	if len(section.Items) != 1 {
		t.Fatalf("expected 1 variable, got %d", len(section.Items))
	}
	variable := section.Items[0]
	if variable.EnvName != "PREFIX_ONE" {
		t.Errorf("expected variable name %q, got %q", "PREFIX_ONE", variable.EnvName)
	}
	if variable.Doc != "First one" {
		t.Errorf("expected variable doc %q, got %q", "First one", variable.Doc)
	}
}

func testRenderer(tmpl template, c renderContext, expectLines ...string) func(*testing.T) {
	return func(t *testing.T) {
		var buf bytes.Buffer
		r := templateRenderer(tmpl)
		err := r(c, &buf)
		if err != nil {
			t.Fatal(err)
		}
		scanner := bufio.NewScanner(&buf)
		var currentLine int
		var logBuilder strings.Builder
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			logBuilder.WriteString(line)
			logBuilder.WriteRune('\n')
			expect := strings.TrimSpace(expectLines[currentLine])
			if line == expect {
				currentLine++
			}
		}
		if err := scanner.Err(); err != nil {
			t.Fatal("error reading output:", err)
		}
		if currentLine != len(expectLines) {
			t.Log(logBuilder.String())
			t.Fatalf("expected line at %d: %q was not found; expected %d lines",
				currentLine, expectLines[currentLine], len(expectLines))
		}
	}
}
