package main

// ComplexConfig is an example configuration structure.
// It contains a few fields with different types of tags.
// It is trying to cover all the possible cases.
//
//go:generate go run ../ -output complex.html -all -format html
//go:generate go run ../ -output complex.md -all
//go:generate go run ../ -output complex.txt -all -format plaintext
type ComplexConfig struct {
	// Secret is a secret value that is read from a file.
	Secret string `env:"SECRET,file"`
	// Password is a password that is read from a file.
	Password string `env:"PASSWORD,file"    envDefault:"/tmp/password"   json:"password"`
	// Certificate is a certificate that is read from a file.
	Certificate string `env:"CERTIFICATE,file,expand" envDefault:"${CERTIFICATE_FILE}"`
	// Key is a secret key.
	SecretKey string `env:"SECRET_KEY,required" json:"secret_key"`
	// SecretVal is a secret value.
	SecretVal string `json:"secret_val" env:"SECRET_VAL,notEmpty"`

	// Hosts is a list of hosts.
	Hosts []string `env:"HOSTS,required" envSeparator:":"`
	// Words is just a list of words.
	Words []string `env:"WORDS,file" envDefault:"one,two,three"`

	Comment string `env:"COMMENT,required" envDefault:"This is a comment."` // Just a comment.
}

// NextConfig is an example configuration structure.
type NextConfig struct { // NextConfig is a configuration structure.
	// Mount is a mount point.
	Mount string `env:"MOUNT,required"`
}
