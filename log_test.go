package main

import "testing"

func TestLog(t *testing.T) {
	flag := debugLogs
	t.Cleanup(func() {
		debugLogs = flag
	})

	if debugLogs {
		t.Fatalf("Expected debugLogs to be false, got %v", debugLogs)
	}
	if l := logger(); l != nullLogger {
		t.Fatalf("Expected nil logger, got %v", l)
	}
	debugLogs = true
	if l := logger(); l != debugLogger {
		t.Fatalf("Expected debug logger, got %v", l)
	}
}
