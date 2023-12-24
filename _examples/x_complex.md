# Environment Variables

## ComplexConfig

ComplexConfig is an example configuration structure.
It contains a few fields with different types of tags.
It is trying to cover all the possible cases.

 - `X_SECRET` (from-file) - Secret is a secret value that is read from a file.
 - `X_PASSWORD` (from-file, default: `/tmp/password`) - Password is a password that is read from a file.
 - `X_CERTIFICATE` (expand, from-file, default: `${CERTIFICATE_FILE}`) - Certificate is a certificate that is read from a file.
 - `X_SECRET_KEY` (**required**) - Key is a secret key.
 - `X_SECRET_VAL` (**required**, non-empty) - SecretVal is a secret value.
 - `X_HOSTS` (separated by `:`, **required**) - Hosts is a list of hosts.
 - `X_WORDS` (comma-separated, from-file, default: `one,two,three`) - Words is just a list of words.
 - `X_COMMENT` (**required**, default: `This is a comment.`) - Just a comment.

## NextConfig

 - `X_MOUNT` (**required**) - Mount is a mount point.

## FieldNames

FieldNames uses field names as env names.

 - `X_QUUX` - Quux is a field with a tag.
