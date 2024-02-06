package main

// Config is the configuration for the application.
type Config struct {
	// Repo is the configuration for the repository.
	Repo struct {
		// Conn is the connection string for the repository.
		Conn string `env:"CONN,notEmpty"`
	} `envPrefix:"REPO_"`
}
