# Examples

There are few target files:
- [`config.go`](./config.go) contains a simple type `Config` with several documented fields. It includes
  `go:generate` directives to call `envdoc` with different params, creating
  documentation in various formats: [`config.md`](./config.md), [`config.html`](./config.html) and [`config.txt`](./config.txt).
- [`complex.go`](./complex.go) features multiple types, covering a broad range of field tags and type comments.
  Additionally, it includes more `go:generate` directives to
  produce documentation in not only different formats but also using various options
  for `envdoc`. For instance, options such as `-no-styles` and `-env-prefix` are utilized here. The outputs include [`complex.md`](./complex.md),
  [`complex.html`](./complex.html) and [`complex.txt`](./complex.txt) for standard markdown, HTML and text documentation, respectively;
  [`x_complex.md`](./x_complex.md) with `-env-prefix` argument;
  and [`complex-nostyle.html`](./complex-nostyle.html) which is HTML documentation without built-in styles.
- [`envprefix.go`](./envprefix.go) showcases a nested config structure with the `envPrefix` tag for a structure field.
  It generates [`envprefix.md`](./envprefix.md).

The examples directory also contains helper script files:
 - `build-examples.sh` - modify any example Go file and regenerate all documentation outputs by executing it via `./build-examples.sh`.
 - `clean.sh` - removes all documentation files.
