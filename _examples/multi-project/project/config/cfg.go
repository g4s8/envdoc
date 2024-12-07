package config

import (
	"lib/logging"
	"project/server"
)

//go:generate go run ../../../../ -files ./config/cfg.go -types * -output ../doc.md -format markdown-table ../ ../../lib/
type Config struct {
	// AppName is the name of the application.
	AppName string `env:"APP_NAME" envDefault:"myapp"`

	// Server config.
	Server server.Config `envPrefix:"SERVER_"`

	// Logging config.
	Log logging.Config `envPrefix:"LOG_"`
}
