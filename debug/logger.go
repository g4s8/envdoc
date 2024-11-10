package debug

import "io"

type Logger interface {
	Logf(format string, args ...interface{})
	Log(args ...interface{})
}

type nopLogger struct{}

func (l *nopLogger) Logf(_ string, _ ...interface{}) {}

func (l *nopLogger) Log(_ ...interface{}) {}

var logger Logger

type Printer interface {
	Debug(io.Writer)
}
