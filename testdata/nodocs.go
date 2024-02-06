package main

type Config struct {
	Repo struct {
		Conn string `env:"CONN,notEmpty"`
	} `envPrefix:"REPO_"`
}
