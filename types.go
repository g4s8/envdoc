package main

type docItemFlags int

const (
	docItemFlagNone     docItemFlags = 0
	docItemFlagRequired docItemFlags = 1 << iota
	docItemFlagExpand
	docItemFlagNonEmpty
	docItemFlagFromFile
)

type docItem struct {
	envName    string // environment variable name
	doc        string // field documentation text
	flags      docItemFlags
	envDefault string
}
