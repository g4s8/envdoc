package testdata

type Type1 struct {
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
	// NotEnv is not an environment variable.
	NotEnv string `json:"not_env"`
	// NoTag is not tagged.
	NoTag string
	// BrokenTag is a tag that is broken.
	BrokenTag string `env:"BROKEN_TAG,required`
}
