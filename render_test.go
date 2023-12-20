package main

import (
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
		strings.Join([]string{
			"# Environment Variables",
			"",
			"- `TEST_ENV` - This is a test environment variable.",
			"- `TEST_ENV2` (comma-separated, default: `default value`) - This is another test environment variable.",
			"- `TEST_ENV3` (**required**, expand, non-empty, from-file) - This is a third test environment variable.",
			"",
		}, "\n")))
	t.Run("plaintext", testRenderer(tmplPlaintext, renderContext{Items: testRenderItems},
		strings.Join([]string{
			"ENVIRONMENT VARIABLES",
			"",
			" * `TEST_ENV` - This is a test environment variable.",
			" * `TEST_ENV2` (comma-separated, default: `default value`) - This is another test environment variable.",
			" * `TEST_ENV3` (required, expand, non-empty, from-file) - This is a third test environment variable.",
			"",
		}, "\n")))
	t.Run("html", testRenderer(tmplHTML, renderContext{Items: testRenderItems},
		strings.Join([]string{
			`<!DOCTYPE html>`,
			`<html lang="en">`,
			`    <head>`,
			`    <meta charset="utf-8">`,
			`    <title>Environment Variables</title>`,
			`    <style>`,
			`    body {`,
			`      font-family: sans-serif;`,
			`    }`,
			`    </style>`,
			`  </head>`,
			`  <body>`,
			`  <h1>Environment Variables</h1>`,
			`  <ul>`,
			`    <li><code>TEST_ENV</code> - This is a test environment variable.</li>`,
			`    <li><code>TEST_ENV2</code> (comma-separated, default: <code>default value</code>) - This is another test environment variable.</li>`,
			`    <li><code>TEST_ENV3</code> (<strong>required</strong>, expand, non-empty, from-file) - This is a third test environment variable.</li>`,
			`  </ul>`,
			`  </body>`,
			`</html>`,
			``,
		}, "\n")))
}

func testRenderer(tmpl template, c renderContext, expect string) func(*testing.T) {
	return func(t *testing.T) {
		var buf bytes.Buffer
		r := templateRenderer(tmpl)
		err := r(c, &buf)
		if err != nil {
			t.Fatal(err)
		}
		if buf.String() != expect {
			t.Errorf("expected %q, got %q", expect, buf.String())
		}
	}
}
