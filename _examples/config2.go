package main

//go:generate go run ../ -output config2.md -type Config2
type Config2 struct {
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
	Hosts []string `env:"HOSTS" envSeparator:":"`
	// Words is just a list of words.
	Words []string `env:"WORDS"`
}
