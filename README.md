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
 * `-template` (path string, *optional*) - Path to a custom template file for rendering the output. It has priority over `-format`.
 * `-title` (string, *optional*, default `Environment Variables`) - Title to be used as the header of the generated file.
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

## Custom Templates
envdoc also supports user-defined templates for more specific layouts provided with the `-template` flag.

Examples can be found [here](_examples/custom-template).

### Template Data
Custom templates are expected to be [Go text templates](https://pkg.go.dev/text/template) and are executed with the following data.

#### Top-Level Data
| Field      | Type              | Description                                                                                             |                                           
|------------|-------------------|---------------------------------------------------------------------------------------------------------|
| `Title`    | `string`          | Value of the `-title` flag. Defaults to `Environment Variables`. Useful as a file header.               |
| `Sections` | `[]renderSection` | A list of structs matched by the `-target` flag.                                                        |
| `Style`    | `bool`            | The opposite of the the `-no-style` flag (hence it defaults to `true`). Useful for toggling CSS styles. |

#### Section: renderSection
Each section represents a struct holding fields that map to environment variables.

| Field    | Type            | Description                                                       |
|----------|-----------------|-------------------------------------------------------------------|
| `Name`   | `string`        | Name of the struct.                                               |
| `Doc`    | `string`        | Description of the struct (parsed from the struct's doc comment). |
| `Items`  | `[]renderItem`  | List of fields within the struct.                                 |

#### Item: renderItem
Each item represents a struct field that maps to an environment variable.

| Field          | Type            | Description                                                                                 |
|----------------|-----------------|---------------------------------------------------------------------------------------------|
| `EnvName`      | `string`        | Name of the environment variable.                                                           |
| `Doc`          | `string`        | Description of the variable (parsed from the field's doc comment).                          |
| `EnvDefault`   | `string`        | Optional default value.                                                                     |
| `EnvSeparator` | `string`        | Character used to separate items in slices and maps.                                        |
| `Required`     | `bool`          | Signifies if the variable must be set.                                                      |
| `NonEmpty`     | `bool`          | Signifies if the variable must not be empty if set. Applies only to `caarlos0`              |
| `Expand`       | `bool`          | Signifies if the value is expandable from environment variables. Applies only to `caarlos0` |
| `FromFile`     | `bool`          | Signifies if the value is read from a file. Applies only to `caarlos0`                      |
| `Children`     | `[]renderItem`  | Nested structs in item. Applies only to `cleanenv`.                                         |

### Functions
Custom templates support the following string functions from the standard library:
- `repeat`
- `split`
- `join`
- `contains`
- `toLower`
- `toUpper`
- `toTitle`
- `replace`
- `hasPrefix`
- `hasSuffix`
- `trimSpace`
- `trimPrefix`
- `trimSuffix`
- `trimLeft`
- `trimRight`

In addition to the standard functions above, the following are supported:
- `strAppend`: `func (arr []string, item string) []string` - Appends `item` to `arr`.
- `strSlice`: `func () []string` - Makes a new empty slice. 
- `list`: `func(args ...any) []any` - Returns the variadic args as a slice.
- `sum`: `func(args ...int) int` - Returns the sum of the variadic args.
- `marshalIndent`: `func(v any) (string, error)` - Marshals the given value into a JSON string. <br>
  Returns an error if the value is not a valid JSON.

## Contributing

If you find any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE.md](/LICENSE.md) file for details.
