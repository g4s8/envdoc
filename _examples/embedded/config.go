package main

import "time"

// Config is the configuration for the application.
//
//go:generate go run ../../ -output doc.md
type Config struct {
	// Start date.
	Start Date `env:"START,notEmpty"`
}

type Time time.Time

// Date is a time.Time wrapper that uses the time.DateOnly layout.
type Date struct {
	Time
}
