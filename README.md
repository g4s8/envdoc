# envdoc

envdoc is a tool for generating documentation for environment variables in Go structs.
It takes comments associated with `env` tags in Go structs and creates a Markdown file with detailed documentation.

This tool is compatible with the
- [caarlos0/env](https://github.com/caarlos0/env) library for parsing environment variables.

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
//go:generate envdoc -output environments.md -type Config
type Config struct {
    // ...
}
```

## Usage

```go
//go:generate envdoc -output <output_file_name> -type <target_type_name> 
```

 * `-output`: Specify the output file name for the generated documentation.
 * `-type`: Specify the target struct type name to generate documentation for.

## Example

Suppose you have the following Go file `config.go`:

```go
package foo

//go:generate envdoc --output env-docs.md --type Config
type Config struct {
  // Port to listen for incoming connections
  Port int `env:"PORT"`
  // Address to serve
  Address string `env:"ADDRESS"`
}
```

And the `go:generate` line above creates documentation in `env-doc.md` file:

```md
# Environment variables

- `PORT` - Port to listen for incoming connections
- `ADDRESS` - Address to serve
```

## Contributing

If you find any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE.md file for details.
