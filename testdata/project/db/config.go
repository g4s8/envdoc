package db

// Config for the database.
type Config struct {
	// Host of the database.
	Host string `env:"HOST,notEmpty"`
	// Port of the database.
	Port int `env:"PORT" envDefault:"5432"`
	// Username of the database.
	Username string `env:"USERNAME,notEmpty"`
	// Password of the database.
	Password string `env:"PASSWORD,notEmpty"`
	// Name of the database.
	Name string `env:"NAME,notEmpty"`
}
