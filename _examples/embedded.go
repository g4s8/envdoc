package main

import "time"

//go:generate go run ../ -output embedded.md
type Config struct {
	// Start date.
	Start Date `env:"START,notEmpty"`
}

// Date is a time.Time wrapper that uses the time.DateOnly layout.
type Date struct {
	time.Time
}
