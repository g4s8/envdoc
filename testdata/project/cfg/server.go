package cfg

type Config struct {
	// Host of the server.
	Host string `env:"HOST,notEmpty"`
	// Port of the server.
	Port int `env:"PORT" envDefault:"8080"`
}
