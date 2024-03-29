package testdata

import "time"

type Config struct {
	// Start date.
	Start Date `env:"START"`
}

// Date is a time.Time wrapper that uses the time.DateOnly layout.
type Date time.Time
