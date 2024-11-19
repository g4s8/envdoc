package testdata

type Config struct {
	// Host is the host of the server.
	Host string `env:"HOST"`

	// Port is the port of the server.
	Port int `env:"PORT"`

	Undocumented string `env:"UNDOCUMENTED"`
}
