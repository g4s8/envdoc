package main

import (
	"fmt"
	"io"
	"strings"
)

type markdownOutput struct {
	out io.Writer

	sb strings.Builder
}

func newMarkdownOutput(out io.Writer) *markdownOutput {
	return &markdownOutput{out: out}
}

func (m *markdownOutput) writeHeader() {
	m.sb.WriteString("# Environment variables\n\n")
}

func (m *markdownOutput) writeItem(item docItem) {
	m.sb.WriteString(fmt.Sprintf("- `%s` ", item.envName))
	var opts []string
	if item.flags&docItemFlagRequired != 0 {
		opts = append(opts, "**required**")
	}
	if item.flags&docItemFlagExpand != 0 {
		opts = append(opts, "expand")
	}
	if item.flags&docItemFlagNonEmpty != 0 {
		opts = append(opts, "not-empty")
	}
	if item.flags&docItemFlagFromFile != 0 {
		opts = append(opts, "from file")
	}
	if item.envDefault != "" {
		opts = append(opts, fmt.Sprintf("default: `%s`", item.envDefault))
	}
	if len(opts) > 0 {
		m.sb.WriteString(fmt.Sprintf("(%s) ", strings.Join(opts, ", ")))
	}
	m.sb.WriteString(fmt.Sprintf("- %s\n", item.doc))
}

func (m *markdownOutput) dump() error {
	if _, err := m.out.Write([]byte(m.sb.String())); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
