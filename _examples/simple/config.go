package main

// Config is an example configuration structure.
// It is used to generate documentation for the configuration
// using the commands below.
//
//go:generate go run ../../ -output doc.txt -format plaintext
//go:generate go run ../../ -output doc.md -format markdown
//go:generate go run ../../ -output doc.html -format html
//go:generate go run ../../ -output doc.env -format dotenv
//go:generate go run ../../ -output doc.env.dist -format dotenv.dist
//go:generate go run ../../ -output doc.json -format json
type Config struct {
	// Hosts name of hosts to listen on.
	Hosts []string `env:"HOST,required", envSeparator:";"`
	// Port to listen on.
	Port int `env:"PORT,notEmpty"`

	// Debug mode enabled.
	Debug bool `env:"DEBUG" envDefault:"false"`

	// Prefix for something.
	Prefix string `env:"PREFIX"`
}
