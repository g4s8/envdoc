package config

import (
	"example.com/db"
	"example.com/server"
)

//go:generate go run ../../../ -dir ../ -files ./config/cfg.go -types * -output ../doc.md -format markdown
type Config struct {
	// AppName is the name of the application.
	AppName string `env:"APP_NAME" envDefault:"myapp"`

	// Server config.
	Server server.Config `envPrefix:"SERVER_"`

	// Database config.
	Database db.Config `envPrefix:"DB_"`

	// Logging config.
	Logging Logging `envPrefix:"LOG_"`
}
