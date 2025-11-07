# envdoc

envdoc is a tool for generating documentation for environment variables in Go structs.
It takes comments associated with `env` tags in Go structs and creates a Markdown, plaintext or HTML
file with detailed documentation.

For `docenv` linter see [docenv/README.md](./docenv/README.md).

<br/>

[![CI](https://github.com/g4s8/envdoc/actions/workflows/go.yml/badge.svg)](https://github.com/g4s8/envdoc/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/g4s8/envdoc.svg)](https://pkg.go.dev/github.com/g4s8/envdoc)
[![codecov](https://codecov.io/gh/g4s8/envdoc/graph/badge.svg?token=sqXWNR755O)](https://codecov.io/gh/g4s8/envdoc)
[![Go Report Card](https://goreportcard.com/badge/github.com/g4s8/envdoc)](https://goreportcard.com/report/github.com/g4s8/envdoc)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  

## Installation

### Go >= 1.24

Add `envdoc` tool and install it:
```bash
go get -tool github.com/g4s8/envdoc@latest
go install tool
```

Add `go:generate`:
```go
//go:generate envdoc -output config.md
type Config struct {
    // ...
}
```

Generate:
```bash
go generate ./...
```

### Before Go 1.24

Run it with `go run` in source file:
```go
//go:generate go run github.com/g4s8/envdoc@latest -output environments.md
type Config struct {
    // ...
}
```

Or download binary to run it:
```bash
go install github.com/g4s8/envdoc@latest
```

And use it in code:

```go
//go:generate envdoc -output environments.md
type Config struct {
    // ...
}
```

## Usage

```go
//go:generate envdoc -output <output_file_name>
```

 * `-dir` (path string, *optional*) - Specify the directory to search for files. Default is the file dir with `go:generate` command.
 * `-files` (glob string, *optional*) - File glob pattern to specify file names to process. Default is the single file with `go:generate`.
 * `-types` (glob string, *optional*) - Type glob pattern for type names to process. If not specified, the next type after `go:generate` is used.
 * `-target` (`enum(caarlos0, cleanenv)` string, optional, default `caarlos0`) - Set env library target.
 * `-output` (path string, **required**) - Output file name for generated documentation.
 * `-format` (`enum(markdown, plaintext, html, dotenv, json)` string, *optional*) - Output format for documentation.  Default is `markdown`.
 * `-no-styles` (`bool`, *optional*) - If true, CSS styles will not be included for `html` format.
 * `-env-prefix` (`string`, *optional*) - Sets additional global prefix for all environment variables.
 * `-tag-name` (string, *optional*, default: `env`) - Use custom tag name instead of `env`.
 * `-tag-default` (string, *optional*, default: `envDefault`) - Use "default" tag name instead of `envDefault`.
 * `-required-if-no-def` (bool, *optional*, default: `false`) - Set attributes as required if no default value is set.
 * `-field-names` (`bool`, *optional*) - Use field names as env names if `env:` tag is not specified.
 * `-debug` (`bool`, *optional*) - Enable debug output.

These params are deprecated and will be removed in the next major release:
 * `-type` - Specify one type to process.
 * `-all` - Process all types in a file.

Both parameters could be replaced with `-types` param:
 - Use `-types=Foo` instead of `-type=Foo`.
 - Use `-types='*'` instead of `-all`.

## Example

Suppose you have the following Go file `config.go`:

```go
package foo

//go:generate envdoc --output env-doc.md
type Config struct {
  // Port to listen for incoming connections
  Port int `env:"PORT,required"`
  // Address to serve
  Address string `env:"ADDRESS" envDefault:"localhost"`
}
```

And the `go:generate` line above creates documentation in `env-doc.md` file:

```md
# Environment Variables

- `PORT` (**required**) - Port to listen for incoming connections
- `ADDRESS` (default: `localhost`) - Address to serve
```

See [_examples](./_examples/) dir for more details.

## Compatibility

This tool is compatible with
- full compatibility: [caarlos0/env](https://github.com/caarlos0/env)
- full compatibility: [ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv)
- partial compatibility: [sethvargo/go-envconfig](https://github.com/sethvargo/go-envconfig)
- partial compatibility: [joeshaw/envdecode](https://github.com/joeshaw/envdecode)

*Let me know about any new lib to check compatibility.*


## Contributing

If you find any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE.md](/LICENSE.md) file for details.
