package main

// Settings is the application settings.
//
//go:generate go run ../ -output envprefix.txt -format plaintext -type Settings
//go:generate go run ../ -output envprefix.md -type Settings
//go:generate go run ../ -output envprefix.html -format html -type Settings
type Settings struct {
	// Database is the database settings
	Database Database `envPrefix:"DB_"`

	// Server is the server settings
	Server ServerConfig `envPrefix:"SERVER_"`

	// Debug is the debug flag
	Debug bool `env:"DEBUG"`
}

// Database is the database settings.
type Database struct {
	// Port is the port to connect to
	Port Int `env:"PORT,required"`
	// Host is the host to connect to
	Host string `env:"HOST,notEmpty" envDefault:"localhost"`
	// User is the user to connect as
	User string `env:"USER"`
	// Password is the password to use
	Password string `env:"PASSWORD"`
	// DisableTLS is the flag to disable TLS
	DisableTLS bool `env:"DISABLE_TLS"`
}

// ServerConfig is the server settings.
type ServerConfig struct {
	// Port is the port to listen on
	Port Int `env:"PORT,required"`

	// Host is the host to listen on
	Host string `env:"HOST,notEmpty" envDefault:"localhost"`

	// Timeout is the timeout settings
	Timeout TimeoutConfig `envPrefix:"TIMEOUT_"`
}

// TimeoutConfig is the timeout settings.
type TimeoutConfig struct {
	// Read is the read timeout
	Read Int `env:"READ" envDefault:"30"`
	// Write is the write timeout
	Write Int `env:"WRITE" envDefault:"30"`
}
