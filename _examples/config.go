package main

//go:generate go run ../ -output config.md -type Config
type Config struct {
	// Host name to listen on.
	Host string `env:"HOST"`
	// Port to listen on.
	Port int `env:"PORT"`

	// Debug mode enabled.
	Debug bool `env:"DEBUG"`
}
