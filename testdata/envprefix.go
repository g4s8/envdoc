package testdata

type Settings struct {
	// Database is the database settings
	Database Database `envPrefix:"DB_"`

	// Debug is the debug flag
	Debug bool `env:"DEBUG"`
}

type Database struct {
	// Port is the port to connect to
	Port Int `env:"PORT,required"`
}
