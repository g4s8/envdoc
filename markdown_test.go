package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestMarkdown(t *testing.T) {
	t.Run("header", func(t *testing.T) {
		var out bytes.Buffer
		mw := newMarkdownOutput(&out)
		mw.writeHeader()
		if err := mw.dump(); err != nil {
			t.Fatal("dump markdown writer", err)
		}
		if out.String() != "# Environment variables\n\n" {
			t.Fatal("unexpected output", out.String())
		}
	})
	t.Run("items", func(t *testing.T) {
		// envName    string // environment variable name
		// doc        string // field documentation text
		// flags      docItemFlags
		// envDefault string
		var out bytes.Buffer
		mw := newMarkdownOutput(&out)
		mw.writeItem(docItem{
			envName: "FOO",
			doc:     "foo",
		})
		mw.writeItem(docItem{
			envName:    "BAR",
			doc:        "bar",
			envDefault: "1",
		})
		mw.writeItem(docItem{
			envName: "BAZ",
			doc:     "baz",
			flags:   docItemFlagRequired,
		})
		mw.writeItem(docItem{
			envName: "QUX",
			doc:     "qux",
			flags:   docItemFlagNonEmpty,
		})
		mw.writeItem(docItem{
			envName: "QUUX",
			doc:     "quux",
			flags:   docItemFlagRequired | docItemFlagNonEmpty,
		})
		mw.writeItem(docItem{
			envName: "CORGE",
			doc:     "corge",
			flags:   docItemFlagRequired | docItemFlagExpand,
		})
		mw.writeItem(docItem{
			envName: "GRAULT",
			doc:     "grault",
			flags:   docItemFlagFromFile | docItemFlagExpand | docItemFlagNonEmpty,
		})
		if err := mw.dump(); err != nil {
			t.Fatal("dump markdown writer", err)
		}
		expectLines := []string{
			"- `FOO` - foo",
			"- `BAR` (default: `1`) - bar",
			"- `BAZ` (**required**) - baz",
			"- `QUX` (not-empty) - qux",
			"- `QUUX` (**required**, not-empty) - quux",
			"- `CORGE` (**required**, expand) - corge",
			"- `GRAULT` (expand, not-empty, from file) - grault",
		}
		scanner := bufio.NewScanner(&out)
		var pos int
		for scanner.Scan() {
			expect := expectLines[pos]
			actual := scanner.Text()
			if actual != expect {
				t.Fatalf("unexpected output at line %d: %q != %q", pos, actual, expect)
			}
			pos++
		}
	})
}
