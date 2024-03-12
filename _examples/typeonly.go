package main

import "time"

//go:generate go run ../ -output typeonly.md -type Config

type Config struct {
	// Some time.
	SomeTime MyTime `env:"SOME_TIME"`
}

type MyTime time.Time

func (t *MyTime) UnmarshalText(text []byte) error {
	tt, err := time.Parse("2006-01-02", string(text))
	*t = MyTime(tt)
	return err
}
