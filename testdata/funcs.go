package testdata

// Test case for #21 where the envdoc panics if target type function presents.

//go:generate envdoc -output test.md --all
type aconfig struct {
	// this is some value
	somevalue string `env:"SOME_VALUE" envDefault:"somevalue"`
}

// when this function is present, go generate panics with "expected type node root child, got nodeField ()".
func someFuncThatTakesInput(configs ...interface{}) {
	// this is some comment
}

func (a *aconfig) someFuncThatTakesInput(configs ...interface{}) {
	// this is some comment
}
