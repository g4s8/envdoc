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

func TestCamelToSnake(t *testing.T) {
	type testCase [2]string
	for _, tc := range []testCase{
		{"Foo", "FOO"},
		{"FooBar", "FOO_BAR"},
		{"FooBarBaz", "FOO_BAR_BAZ"},
		{"fooBar", "FOO_BAR"},
		{"fooBarBaz", "FOO_BAR_BAZ"},
		{"", ""},
		{"Foo_", "FOO_"},
		{"Foo_Bar", "FOO_BAR"},
		{"Foo_bar", "FOO_BAR"},
		{"Foo_Bar_Baz", "FOO_BAR_BAZ"},
	} {
		if got := camelToSnake(tc[0]); got != tc[1] {
			t.Fatalf("expected camelToSnake(%q) to be %q, got %q", tc[0], tc[1], got)
		}
	}
}
