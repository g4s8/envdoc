package testdata

type Foo struct {
	// One is a one.
	One string `env:"ONE"`
	// Two is a two.
	Two string `env:"TWO"`
}

// Bar is a bar.
type Bar struct {
	// Three is a three.
	Three string `env:"THREE"`
	// Four is a four.
	Four string `env:"FOUR"`
}
