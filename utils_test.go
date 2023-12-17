package main

import (
	"errors"
	"testing"
)

type mockCloser struct {
	err    error
	closed bool
}

func (m *mockCloser) Close() error {
	m.closed = true
	return m.err
}

func TestCloseWith(t *testing.T) {
	var (
		targetErr  = errors.New("test")
		handlerErr error
		closer     = &mockCloser{err: targetErr}
		handler    = func(err error) {
			handlerErr = err
		}
	)
	closeWith(closer, handler)
	if !closer.closed {
		t.Fatal("expected closer to be closed")
	}
	if handlerErr != targetErr {
		t.Fatalf("expected handler to be called with %q, got %q", targetErr, handlerErr)
	}
}
