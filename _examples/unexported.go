package main

//go:generate go run ../ -output unexported.md
type appconfig struct {
	// Port the application will listen on inside the container
	Port int `env:"PORT" envDefault:"8080"`
	// some more stuff I omitted here
}
