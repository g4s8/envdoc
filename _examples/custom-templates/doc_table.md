# Configuration

## Config

Config is an example configuration structure.
It is used to generate documentation from custom templates.

| Name | Description | Default | Attributes |
|------|-------------|---------|------------|
| `SECRET` | Secret is a secret value that is read from a file. |  | `From File` |
| `PASSWORD` | Password is a password that is read from a file. | `/tmp/password` | `From File` |
| `CERTIFICATE` | Certificate is a certificate that is read from a file. | `${CERTIFICATE_FILE}` | `Expandable`, `From File` |
| `SECRET_KEY` | Key is a secret key. |  | `REQUIRED` |
| `SECRET_VAL` | SecretVal is a secret value. |  | `REQUIRED`, `Not Empty` |
| `HOSTS` | Hosts is a list of hosts. |  | `REQUIRED`, `Separated by :` |
| `WORDS` | Words is just a list of words. | `one,two,three` | `From File`, `Separated by ,` |
| `COMMENT` | Just a comment. | `This is a comment.` | `REQUIRED` |
| `ALLOW_METHODS` | AllowMethods is a list of allowed methods. | `GET, POST, PUT, PATCH, DELETE, OPTIONS` |  |
| `ANON_USER` | User is a user name. |  | `REQUIRED` |
| `ANON_PASS` | Pass is a password. |  | `REQUIRED` |

## NextConfig

NextConfig is a configuration structure to generate multiple doc sections.

| Name | Description | Default | Attributes |
|------|-------------|---------|------------|
| `MOUNT` | Mount is a mount point. |  | `REQUIRED` |
