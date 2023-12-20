package main

import "testing"

type formatExpectation struct {
	before string
	items  []string
	after  string
}

func TestFormats(t *testing.T) {
	items := []docItem{
		{
			envName: "FOO",
			doc:     "foo",
		},
		{
			envName:    "BAR",
			doc:        "bar",
			envDefault: "1",
			separator:  ";",
			flags:      docItemFlagRequired | docItemFlagExpand | docItemFlagNonEmpty | docItemFlagFromFile,
		},
		{
			envName:   "BAZ",
			doc:       "baz",
			separator: ",",
		},
	}
	t.Run("markdown", formatTester(fmtMD, items, formatExpectation{
		before: "# Environment Variables\n\n",
		items: []string{
			"- `FOO` - foo\n",
			"- `BAR` (separated by `;`, **required**, expand, not-empty, from file, default: `1`) - bar\n",
			"- `BAZ` (comma-separated) - baz\n",
		},
		after: "",
	}))
	t.Run("plaintext", formatTester(fmtPlain, items, formatExpectation{
		before: "ENVIRONMENT VARIABLES\n\n",
		items: []string{
			" * `FOO` - foo\n",
			" * `BAR` (separated by `;`, required, expand, not-empty, from file, default: `1`) - bar\n",
			" * `BAZ` (comma-separated) - baz\n",
		},
		after: "",
	}))
	t.Run("html", formatTester(fmtHTML, items, formatExpectation{
		before: `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Environment Variables</title>
		<style>
		body {
			font-family: sans-serif;
		}
		</style>
	</head>
	<body>
		<h1>Environment Variables</h1>
		<ul>
`,
		items: []string{
			"<li><code>FOO</code> - foo</li>\n",
			"<li><code>BAR</code> (separated by <code>;</code>, <strong>required</strong>, expand, not-empty, from file, default: <code>1</code>) - bar</li>\n",
			"<li><code>BAZ</code> (comma-separated) - baz</li>\n",
		},
		after: `</ul>
	</body>
</html>
`,
	}))
}

func formatTester(f itemFormat, items []docItem, expect formatExpectation) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("before", func(t *testing.T) {
			actual := f.beforeItems()
			if actual != expect.before {
				t.Fatalf("unexpected output: %q != %q", actual, expect.before)
			}
		})
		t.Run("items", func(t *testing.T) {
			for i, item := range items {
				actual := f.formatItem(item)
				if actual != expect.items[i] {
					t.Fatalf("unexpected output at line %d: %q != %q", i, actual, expect.items[i])
				}
			}
		})
		t.Run("after", func(t *testing.T) {
			actual := f.afterItems()
			if actual != expect.after {
				t.Fatalf("unexpected output: %q != %q", actual, expect.after)
			}
		})
	}
}
