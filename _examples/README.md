# Examples

There are two target files: `config.go` and `complex.go`:
- `config.go` has a simple type `Config` with a few documented fields. It has
  `go:generate` directives to call `envdoc` with different params to build
  documentation in different formats.
- `complex.go` has multiple types, and this file tries to cover all possible cases
  of field tags, type comments, etc. Also, it has more `go:generate` directives to
  produce documentation not only in different formats but it uses different options
  for `envdoc`, e.g. there you can find `-no-styles` and `-env-prefix` options.
