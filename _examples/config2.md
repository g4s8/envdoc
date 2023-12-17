# Environment Variables

- `SECRET` (from file) - Secret is a secret value that is read from a file.
- `PASSWORD` (from file, default: `/tmp/password`) - Password is a password that is read from a file.
- `CERTIFICATE` (expand, from file, default: `${CERTIFICATE_FILE}`) - Certificate is a certificate that is read from a file.
- `SECRET_KEY` (**required**) - Key is a secret key.
- `SECRET_VAL` (**required**, not-empty) - SecretVal is a secret value.
