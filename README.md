# envdoc

envdoc is a tool for generating documentation for environment variables in Go structs.
It takes comments associated with `env` tags in Go structs and creates a Markdown, plaintext or HTML
file with detailed documentation.


<br/>

[![CI](https://github.com/g4s8/envdoc/actions/workflows/go.yml/badge.svg)](https://github.com/g4s8/envdoc/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/g4s8/envdoc.svg)](https://pkg.go.dev/github.com/g4s8/envdoc)
[![codecov](https://codecov.io/gh/g4s8/envdoc/graph/badge.svg?token=sqXWNR755O)](https://codecov.io/gh/g4s8/envdoc)
[![Go Report Card](https://goreportcard.com/badge/github.com/g4s8/envdoc)](https://goreportcard.com/report/github.com/g4s8/envdoc)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  

## Installation

Run it with `go run` in source file:
```go
//go:generate go run github.com/g4s8/envdoc@latest -output environments.md -type Config
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
//go:generate envdoc -output <output_file_name> -type <target_type_name> 
```

 * `-output` (**required**) - Specify the output file name for the generated documentation.
 * `-type`: Specify the target struct type name to generate documentation for.
 If ommited, the next type after `go:generate` comment will be used.
 * `-format` (default: `markdown`) - Set output format type, either `markdown`,
 `plaintext`, `html`, or `dotenv`.
 * `-all` - Generate documentation for all types in the file.
 * `-env-prefix` - Environment variable prefix.
 * `-no-styles` - Disable built-int CSS styles for HTML format.
 * `-field-names` - Use field names instead of struct tags for variable names
   if tags are not set.

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
- partial compatibility: [sethvargo/go-envconfig](https://github.com/sethvargo/go-envconfig)
- partial compatibility: [joeshaw/envdecode](https://github.com/joeshaw/envdecode)

*Let me know about any new lib to check compatibility.*


## Contributing

If you find any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE.md](/LICENSE.md) file for details.
