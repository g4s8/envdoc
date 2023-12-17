package main

//go:generate go run ../ -output config.html -format html
type Config struct {
	// Host name to listen on.
	Host string `env:"HOST,required"`
	// Port to listen on.
	Port int `env:"PORT,notEmpty"`

	// Debug mode enabled.
	Debug bool `env:"DEBUG" envDefault:"false"`
}
