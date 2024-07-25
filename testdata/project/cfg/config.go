package cfg

import "github.com/smallstep/certificates/db"

// Config for the application.
type Config struct {
	// Environment of the application.
	Environment string `env:"ENVIRONMENT,notEmpty" envDefault:"development"`

	ServerConfig `envPrefix:"SERVER_"`

	// Database config.
	Database db.Config `envPrefix:"DB_"`
}
