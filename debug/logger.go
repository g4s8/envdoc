package debug

import "io"

type Logger interface {
	Logf(format string, args ...any)
	Log(args ...any)
}

type nopLogger struct{}

func (l *nopLogger) Logf(_ string, _ ...any) {}

func (l *nopLogger) Log(_ ...any) {}

var logger Logger

type Printer interface {
	Debug(io.Writer)
}
