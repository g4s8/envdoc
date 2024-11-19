package testdata

type Config struct {
	// Host is the host of the server.
	Host string `foo:"HOST"`

	// Port is the port of the server.
	Port int `foo:"PORT"`

	Undocumented string `foo:"UNDOCUMENTED"`
}
