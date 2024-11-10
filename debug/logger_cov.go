//go:build coverage

package debug

import (
	"io"
	"testing"
)

func NewLogger(out io.Writer) Logger {
	return &nopLogger{}
}

func SetLogger() {
}

func SetTestLogger(t *testing.T) {
}

func Logf(format string, args ...interface{}) {
}

func Log(args ...interface{}) {
}

func PrintDebug(p Printer) {
}
