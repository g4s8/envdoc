//go:build !coverage

package debug

import (
	"fmt"
	"io"
	"os"
	"testing"
)

type ioLogger struct {
	out io.Writer
}

func NewLogger(out io.Writer) Logger {
	return &ioLogger{out: out}
}

func (l *ioLogger) Logf(format string, args ...interface{}) {
	fmt.Fprintf(l.out, format, args...)
}

func (l *ioLogger) Log(args ...interface{}) {
	fmt.Fprint(l.out, args...)
}

type testLogger struct {
	t *testing.T
}

func NewTestLogger(t *testing.T) Logger {
	return &testLogger{t: t}
}

func (l *testLogger) Logf(format string, args ...interface{}) {
	l.t.Helper()
	l.t.Logf(format, args...)
}

func (l *testLogger) Log(args ...interface{}) {
	l.t.Helper()
	l.t.Log(args...)
}

func SetLogger() {
	if !Config.Enabled {
		logger = &nopLogger{}
		return
	}
	logger = NewLogger(os.Stdout)
}

func SetTestLogger(t *testing.T) {
	if logger == nil {
		SetLogger()
	}
	currentLogger := logger
	t.Cleanup(func() {
		logger = currentLogger
	})
	logger = NewTestLogger(t)
}

func Logf(format string, args ...interface{}) {
	if logger == nil {
		SetLogger()
	}
	logger.Logf(format, args...)
}

func Log(args ...interface{}) {
	if logger == nil {
		SetLogger()
	}
	logger.Log(args...)
}

func PrintDebug(p Printer) {
	if Config.Enabled {
		p.Debug(os.Stdout)
	}
}
