# Environment Variables

## ComplexConfig

ComplexConfig is an example configuration structure.
It contains a few fields with different types of tags.
It is trying to cover all the possible cases.

 - `SECRET` (from-file) - Secret is a secret value that is read from a file.
 - `PASSWORD` (from-file, default: `/tmp/password`) - Password is a password that is read from a file.
 - `CERTIFICATE` (expand, from-file, default: `${CERTIFICATE_FILE}`) - Certificate is a certificate that is read from a file.
 - `SECRET_KEY` (**required**) - Key is a secret key.
 - `SECRET_VAL` (**required**, non-empty) - SecretVal is a secret value.
 - `HOSTS` (separated by `:`, **required**) - Hosts is a list of hosts.
 - `WORDS` (comma-separated, from-file, default: `one,two,three`) - Words is just a list of words.
 - `COMMENT` (**required**, default: `This is a comment.`) - Just a comment.
 - Anon is an anonymous structure.
   - `ANON_USER` (**required**) - User is a user name.
   - `ANON_PASS` (**required**) - Pass is a password.

## NextConfig

 - `MOUNT` (**required**) - Mount is a mount point.

## FieldNames

FieldNames uses field names as env names.

 - `QUUX` - Quux is a field with a tag.

## Config

 - `HOST` (separated by `;`, **required**) - Hosts name of hosts to listen on.
 - `PORT` (**required**, non-empty) - Port to listen on.
 - `DEBUG` (default: `false`) - Debug mode enabled.
 - `PREFIX` - Prefix for something.

## Config

 - `START` (**required**, non-empty) - Start date.

## Date

Date is a time.Time wrapper that uses the time.DateOnly layout.


## Settings

Settings is the application settings.

 - Database is the database settings
   - `DB_PORT` (**required**) - Port is the port to connect to
   - `DB_HOST` (**required**, non-empty, default: `localhost`) - Host is the host to connect to
   - `DB_USER` - User is the user to connect as
   - `DB_PASSWORD` - Password is the password to use
   - `DB_DISABLE_TLS` - DisableTLS is the flag to disable TLS
 - Server is the server settings
   - `SERVER_PORT` (**required**) - Port is the port to listen on
   - `SERVER_HOST` (**required**, non-empty, default: `localhost`) - Host is the host to listen on
   - Timeout is the timeout settings
     - `SERVER_TIMEOUT_READ` (default: `30`) - Read is the read timeout
     - `SERVER_TIMEOUT_WRITE` (default: `30`) - Write is the write timeout
 - `DEBUG` - Debug is the debug flag

## Database

Database is the database settings.

 - `PORT` (**required**) - Port is the port to connect to
 - `HOST` (**required**, non-empty, default: `localhost`) - Host is the host to connect to
 - `USER` - User is the user to connect as
 - `PASSWORD` - Password is the password to use
 - `DISABLE_TLS` - DisableTLS is the flag to disable TLS

## ServerConfig

ServerConfig is the server settings.

 - `PORT` (**required**) - Port is the port to listen on
 - `HOST` (**required**, non-empty, default: `localhost`) - Host is the host to listen on
 - Timeout is the timeout settings
   - `TIMEOUT_READ` (default: `30`) - Read is the read timeout
   - `TIMEOUT_WRITE` (default: `30`) - Write is the write timeout

## TimeoutConfig

TimeoutConfig is the timeout settings.

 - `READ` (default: `30`) - Read is the read timeout
 - `WRITE` (default: `30`) - Write is the write timeout

## Config

 - `START` - Start date.

## Date

Date is a time.Time wrapper that uses the time.DateOnly layout.


## appconfig

 - `PORT` (default: `8080`) - Port the application will listen on inside the container

