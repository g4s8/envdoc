package main

// Struct for tag customization.
//
//go:generate go run ../../ -output ./doc.md -tag-name xenv -tag-default xdef -required-if-no-def
type CustomTagsConfig struct {
	// Host is the host name.
	Host string `xenv:"host" xdef:"localhost"`
	// NoDef is the no default value.
	NoDef string `xenv:"no_def"`
}
