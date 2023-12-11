package main

//go:generate go run ../ -output config2.md -type Config2
type Config2 struct {
	// Secret is a secret value that is read from a file.
	Secret string `env:"SECRET,file"`
	// Password is a password that is read from a file.
	Password string `env:"PASSWORD,file"           envDefault:"/tmp/password"   json:"password"`
	// Certificate is a certificate that is read from a file.
	Certificate string `env:"CERTIFICATE,file,expand" envDefault:"${CERTIFICATE_FILE}"`
	// Key is a secret key.
	SecretKey string `env:"SECRET_KEY,required" json:"secret_key"`
	// SecretVal is a secret value.
	SecretVal string `json:"secret_val" env:"SECRET_VAL,notEmpty"`
	// JustAMess is a mess of env options.
	JustAMess string `fooo:"1" env:"JUST_A_MESS,required,notEmpty,file,expand" json:"just_a_mess" envDefault:"${JUST_A_MESS_FILE}" bar:"2"`
}
