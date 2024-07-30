package server

// TimeoutConfig holds the configuration for the timeouts of the server.
type TimeoutConfig struct {
	// ReadTimeout of the server.
	ReadTimeout string `env:"READ,required"`
	// WriteTimeout of the server.
	WriteTimeout string `env:"WRITE,required"`
}
