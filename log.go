package main

import (
	"log"
	"os"
)

type nullWriter struct{}

func (nullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

var (
	nullLogger  = log.New(nullWriter{}, "", 0)
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
)

var debugLogs = false

func logger() *log.Logger {
	if debugLogs {
		return debugLogger
	}
	return nullLogger
}
