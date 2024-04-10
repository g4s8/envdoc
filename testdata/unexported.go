package testdata

type appconfig struct {
	// Port the application will listen on inside the container
	Port int `env:"PORT" envDefault:"8080"`
}
