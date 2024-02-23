package testdata

// FieldNames uses field names as env names.
type FieldNames struct {
	// Foo is a single field.
	Foo string
	// Bar and Baz are two fields.
	Bar, Baz string
	// Quux is a field with a tag.
	Quux string `env:"QUUX"`
	// FooBar is a field with a default value.
	FooBar string `envDefault:"quuux"`
	// Required is a required field.
	Required string `env:",required"`
}
