# docenv - linter for environment documentation

The linter check that all environment variable fields with `env` tag are documented.

## Install linter

```
go install github.com/g4s8/envdoc/docenv@latest
```

## Example

The struct with undocumented fields:
```go
type Config struct {
	Hosts []string `env:"HOST,required", envSeparator:";"`
	Port  int      `env:"PORT,notEmpty"`
}
```

Run the linter:
```bash
$ go vet -vettool=$(which docenv) ./config.go
config.go:12:2: field `Hosts` with `env` tag should have a documentation comment
config.go:13:2: field `Port` with `env` tag should have a documentation comment
```

## Usage

Flags:
 - `env-name` sets custom env tag name (default `env`)

```
go vet go vet -vettool=$(which docenv) -docenv.env-name=foo ./config.go
```
