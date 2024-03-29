package main

import "time"

//go:generate go run ../ -output typedef.md
type Config struct {
	// Start date.
	Start Date `env:"START"`
}

// Date is a time.Time wrapper that uses the time.DateOnly layout.
type Date time.Time
