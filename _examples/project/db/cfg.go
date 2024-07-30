package db

// Config holds the configuration for the database.
type Config struct {
	// Host of the database.
	Host string `env:"HOST,required"`
	// Port of the database.
	Port string `env:"PORT,required"`
	// User of the database.
	User string `env:"USER" envDefault:"user"`
	// Password of the database.
	Password string `env:"PASSWORD,nonempty"`

	SslConfig `envPrefix:"SSL_"`
}

// SslConfig holds the configuration for the SSL of the database.
type SslConfig struct {
	// SslMode of the database.
	SslMode string `env:"MODE" envDefault:"disable"`
	// SslCert of the database.
	SslCert string `env:"CERT"`
	// SslKey of the database.
	SslKey string `env:"KEY"`
}
