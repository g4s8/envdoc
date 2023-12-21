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

func TestRender(t *testing.T) {
	t.Run("markdown", testRenderer(tmplMarkdown, renderContext{Items: testRenderItems},
		"# Environment Variables",
		"- `TEST_ENV` - This is a test environment variable.",
		"- `TEST_ENV2` (comma-separated, default: `default value`) - This is another test environment variable.",
		"- `TEST_ENV3` (**required**, expand, non-empty, from-file) - This is a third test environment variable."))
	t.Run("plaintext", testRenderer(tmplPlaintext, renderContext{Items: testRenderItems},
		"ENVIRONMENT VARIABLES",
		" * `TEST_ENV` - This is a test environment variable.",
		" * `TEST_ENV2` (comma-separated, default: `default value`) - This is another test environment variable.",
		" * `TEST_ENV3` (required, expand, non-empty, from-file) - This is a third test environment variable."))
	t.Run("html", testRenderer(tmplHTML, renderContext{Items: testRenderItems},
		`<!DOCTYPE html>`,
		`<html lang="en">`,
		`<head>`,
		`<meta charset="utf-8">`,
		`<title>Environment Variables</title>`,
		`</head>`,
		`<section>`,
		`<article>`,
		`<h1>Environment Variables</h1>`,
		`<ul>`,
		`<li><code>TEST_ENV</code> - This is a test environment variable.</li>`,
		`<li><code>TEST_ENV2</code> (comma-separated, default: <code>default value</code>) - This is another test environment variable.</li>`,
		`<li><code>TEST_ENV3</code> (<strong>required</strong>, expand, non-empty, from-file) - This is a third test environment variable.</li>`,
		`</ul>`,
		`</article>`,
		`</section>`,
		`</body>`,
		`</html>`))
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
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			expect := strings.TrimSpace(expectLines[currentLine])
			if line == expect {
				currentLine++
			}
		}
		if err := scanner.Err(); err != nil {
			t.Fatal("error reading output:", err)
		}
		if currentLine != len(expectLines) {
			t.Log("output:")
			t.Log(buf.String())
			t.Fatalf("expected %d lines, got %d", len(expectLines), currentLine)
		}
	}
}
