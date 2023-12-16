package testdata

//go:generate STUB
type Type1 struct {
	// Foo stub
	Foo int `env:"FOO"`
}

type Type2 struct {
	// Baz stub
	Baz int `env:"BAZ"`
}
