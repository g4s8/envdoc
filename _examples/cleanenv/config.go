package main

// Config is an example configuration structure.
// It is used to generate documentation for the configuration
// using the commands below.
//
//go:generate go run ../../ -output doc.txt -target cleanenv -format plaintext
type Config struct {
	// Hosts name of hosts to listen on.
	Hosts []string `env:"HOST" env-required:"true" env-separator:";"`
	// Port to listen on.
	Port int `env:"PORT"`

	// Debug mode enabled.
	Debug bool `env:"DEBUG" env-default:"false"`

	// Location of the server.
	Location string `env:"LOCATION" env-default:"city,country"`

	// Timeouts configuration.
	Timeouts struct {
		// Read timeout.
		Read int `env:"READ" env-default:"10"`
		// Write timeout.
		Write int `env:"WRITE" env-default:"10"`
	} `env-prefix:"TIMEOUT_"`
}
