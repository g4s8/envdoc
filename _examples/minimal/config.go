package minimal

// OAuthConfig holds configuration for OAuth clients and auth redirects.
//
//go:generate go run ../../ -output doc.txt -format plaintext
//go:generate go run ../../ -output doc.md -format markdown
//go:generate go run ../../ -output doc.html -format html
//go:generate go run ../../ -output doc.env -format dotenv
//go:generate go run ../../ -output doc.json -format json
type OAuthConfig struct {
	App struct {
		ID     string   `env:"ID,notEmpty"`
		Secret string   `env:"SECRET,notEmpty" envDefault:"changeme"`
		Scopes []string `env:"SCOPES" envSeparator:" "`
	} `envPrefix:"APP_"`

	Auth struct {
		Redirect struct {
			External      string `env:"EXTERNAL_URL" envDefault:"http://localhost/"`
			InternalRoute string `env:"INTERNAL_ROUTE" envDefault:""`
		} `envPrefix:"REDIRECT_"`
	} `envPrefix:"AUTH_"`

	GitHubIssueTesting struct {
		Foo string `env:"FOO" envDefault:""`
		Bar string `env:"BAR" envDefault:"abc"`
	} `envPrefix:"TESTING_"`
}
