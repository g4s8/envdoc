package testutils

import "testing"

func AssertFatal(t *testing.T, ok bool, format string, args ...interface{}) {
	if !ok {
		t.Fatalf(format, args...)
	}
}

func AssertError(t *testing.T, ok bool, format string, args ...interface{}) {
	if !ok {
		t.Errorf(format, args...)
	}
}
