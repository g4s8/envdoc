package main

import "sync/atomic"

type Foo struct {
	X atomic.Bool
	T bool
	F func(int) string
}
