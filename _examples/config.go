package main

// Config is an example configuration structure.
// It is used to generate documentation for the configuration
// using the commands below.
//
//go:generate go run ../ -output config.txt -format plaintext
//go:generate go run ../ -output config.md -format markdown
//go:generate go run ../ -output config.html -format html
type Config struct {
	// Host name to listen on.
	Host string `env:"HOST,required"`
	// Port to listen on.
	Port int `env:"PORT,notEmpty"`

	// Debug mode enabled.
	Debug bool `env:"DEBUG" envDefault:"false"`
}
