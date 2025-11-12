package main

// Config is an example configuration structure.
// It is used to generate documentation from custom templates.
//
//go:generate go run ../../ -output doc_table.md -template mdtable.tmpl -types *
//go:generate go run ../../ -output doc_tree.txt -target cleanenv -template texttree.tmpl -types *
//go:generate go run ../../ -output doc_table_styled.html -template htmltable.tmpl -types *
//go:generate go run ../../ -output doc_table_plain.html -template htmltable.tmpl -types * -no-styles true
type Config struct {
	// Secret is a secret value that is read from a file.
	Secret string `env:"SECRET,file"`
	// Password is a password that is read from a file.
	Password string `env:"PASSWORD,file" envDefault:"/tmp/password" env-default:"/tmp/password" json:"password"`
	// Certificate is a certificate that is read from a file.
	Certificate string `env:"CERTIFICATE,file,expand" envDefault:"${CERTIFICATE_FILE}" env-default:"${CERTIFICATE_FILE}"`
	// Key is a secret key.
	SecretKey string `env:"SECRET_KEY,required" env-required:"true" json:"secret_key"`
	// SecretVal is a secret value.
	SecretVal string `json:"secret_val" env:"SECRET_VAL,notEmpty"`

	// Hosts is a list of hosts.
	Hosts []string `env:"HOSTS,required" env-required:"true" envSeparator:":" env-separator:":"`

	// Words is just a list of words.
	Words []string `env:"WORDS,file" envDefault:"one,two,three" env-default:"one,two,three"`

	Comment string `env:"COMMENT,required" env-required:"true" envDefault:"This is a comment." env-default:"This is a comment."` // Just a comment.

	// AllowMethods is a list of allowed methods.
	AllowMethods string `env:"ALLOW_METHODS" envDefault:"GET, POST, PUT, PATCH, DELETE, OPTIONS" env-default:"GET, POST, PUT, PATCH, DELETE, OPTIONS"`

	// Anon is an anonymous structure.
	Anon struct {
		// User is a user name.
		User string `env:"USER,required" env-required:"true"`
		// Pass is a password.
		Pass string `env:"PASS,required" env-required:"true"`
	} `envPrefix:"ANON_"`
}

// NextConfig is a configuration structure to generate multiple doc sections.
type NextConfig struct {
	// Mount is a mount point.
	Mount string `env:"MOUNT,required" env-required:"true"`
}
