package server

type Config struct {
	// Host of the server.
	Host string `env:"HOST,required"`
	// Port of the server.
	Port string `env:"PORT,required"`
	// Timeout of the server.
	Timeout TimeoutConfig `envPrefix:"TIMEOUT_"`
}
