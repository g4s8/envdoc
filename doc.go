/*
envdoc is a tool to generate documentation for environment variables
from a Go source file. It is intended to be used as a go generate
directive.

For example, given the following Go type with struct tags and
a `go:generate` directive:

	//go:generate go run github.com/g4s8/envdoc@latest -output config.md
	type Config struct {
		// Host name to listen on.
		Host string `env:"HOST,required"`
		// Port to listen on.
		Port int `env:"PORT,notEmpty"`

		// Debug mode enabled.
		Debug bool `env:"DEBUG" envDefault:"false"`
	}

Running go generate will generate the following Markdown file:

	# Environment variables

	- `HOST` (**required**) - Host name to listen on.
	- `PORT` (**required**, not-empty) - Port to listen on.
	- `DEBUG` (default: `false`) - Debug mode enabled.

By default envdoc generates documentation in Markdown format, but it
can also generate plaintext or HTML.

Options:
  - `-output` - Output file name.
  - `-type` - Type name to generate documentation for. Defaults for
    the next type after `go:generate` directive.
  - `-format` (default: `markdown`) - Set output format type, either `markdown`,
    `plaintext` or `html`.
  - `-all` - Generate documentation for all types in the file.
  - `-env-prefix` - Environment variable prefix.
*/
package main
