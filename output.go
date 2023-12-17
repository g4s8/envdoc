package main

import (
	"fmt"
	"io"
)

type itemFormat interface {
	formatItem(item docItem) string
	beforeItems() string
	afterItems() string
}

type docOutput struct {
	out    io.Writer
	format itemFormat
}

func newDocOutput(out io.Writer, format itemFormat) *docOutput {
	return &docOutput{out: out, format: format}
}

func (m *docOutput) begin() {
	fmt.Fprint(m.out, m.format.beforeItems())
}

func (m *docOutput) writeItem(item docItem) {
	fmt.Fprint(m.out, m.format.formatItem(item))
}

func (m *docOutput) end() {
	fmt.Fprint(m.out, m.format.afterItems())
}
