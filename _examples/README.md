# Examples

There are few target files:
- [`config.go`](./config.go) has a simple type `Config` with a few documented fields. It has
  `go:generate` directives to call `envdoc` with different params to build
  documentation in different formats: [`config.md`](./config.md), [`config.html`](./config.html) and [`config.txt`](./config.txt).
- [`complex.go`](./complex.go) has multiple types, and this file tries to cover all possible cases
  of field tags and type comments. Also, it has more `go:generate` directives to
  produce documentation not only in different formats but it uses different options
  for `envdoc`, e.g. there you can find `-no-styles` and `-env-prefix` options: [`complex.md`](./complex.md), [`complex.html`](./complex.html),
  [`complex.txt`](./complex.txt) for default markdown, HTML and text documentation files; [`x_complex.md`](./x_complex.md) with
  envprefix parameter; [`complex-nostyle.html`](./complex-nostyle.html) is HTML documentation without built-in styles.
- [`envprefix.go`](./envprefix.go) has nested config structure with `envPrefix` tag for structure field. It generates [`envprefix.md`](./envprefix.md).

The examples dir has helper script files as well:
 - `build-examples.sh` - you can modify any example go file and regenerate all documentation outputs by running it via `./build-examples.sh`.
 - `clean.sh` - removes all documentation files.
