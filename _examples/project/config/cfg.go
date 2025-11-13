package config

import (
	"example.com/db"
	srv "github.com/docker/docker/api/server"
)

//go:generate go run ../../.. -dir ../ -files ./config/cfg.go -types * -output ../config.md -format markdown -title Configuration
type Config struct {
	// AppName is the name of the application.
	AppName string `env:"APP_NAME" envDefault:"myapp"`

	// Server config.
	Server srv.Config `envPrefix:"SERVER_"`

	// Database config.
	Database db.Config `envPrefix:"DB_"`

	// Logging config.
	Logging Logging `envPrefix:"LOG_"`
}
