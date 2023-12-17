package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

type testFormat int

func (testFormat) formatItem(item docItem) string {
	return fmt.Sprintf("%s|%s\n", item.envName, item.doc)
}

func (testFormat) beforeItems() string {
	return "#BEFORE\n"
}

func (testFormat) afterItems() string {
	return "#AFTER\n"
}

func TestMarkdown(t *testing.T) {
	var out bytes.Buffer
	mw := newDocOutput(&out, testFormat(0))
	mw.begin()
	mw.writeItem(docItem{
		envName: "FOO",
		doc:     "foo",
	})
	mw.writeItem(docItem{
		envName: "BAR",
		doc:     "bar",
	})
	mw.end()
	expectLines := []string{
		"#BEFORE",
		"FOO|foo",
		"BAR|bar",
		"#AFTER",
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
}
