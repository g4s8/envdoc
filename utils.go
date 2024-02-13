package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"unicode"
)

func closeWith(closer io.Closer, handler func(error)) {
	err := closer.Close()
	if err != nil {
		handler(err)
	}
}

func camelToSnake(s string) string {
	const underscore = '_'
	var result strings.Builder
	result.Grow(len(s) + 5)

	var prev rune
	for i, r := range s {
		if i > 0 && prev != underscore && r != underscore && unicode.IsUpper(r) {
			result.WriteRune(underscore)
		}
		result.WriteRune(unicode.ToUpper(r))
		prev = r
	}

	return result.String()
}

func fastRandString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seed := rand.Intn(len(letters)*len(letters)) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[(seed+i)%len(letters)]
	}
	return string(b)
}

type (
	envDocItemBuilderOp func(*envDocItemsBuilder)
	envDocItemsBuilder  struct {
		envPrefix string
		names     []string
		doc       string
		opts      EnvVarOptions
		children  []*EnvDocItem
	}
)

func withEnvDocItemEnvPrefix(envPrefix string) envDocItemBuilderOp {
	return func(b *envDocItemsBuilder) {
		b.envPrefix = envPrefix
	}
}

func withEnvDocItemDoc(doc string) envDocItemBuilderOp {
	return func(b *envDocItemsBuilder) {
		b.doc = doc
	}
}

func withEnvDocItemOpts(opts EnvVarOptions) envDocItemBuilderOp {
	return func(b *envDocItemsBuilder) {
		b.opts = opts
	}
}

func withEnvDocItemAddChildren(children []*EnvDocItem) envDocItemBuilderOp {
	return func(b *envDocItemsBuilder) {
		b.children = append(b.children, children...)
	}
}

func withEnvDocItemNames(names ...string) envDocItemBuilderOp {
	return func(b *envDocItemsBuilder) {
		b.names = names
	}
}

var withEnvDocEmptyNames = withEnvDocItemNames("")

func (b *envDocItemsBuilder) apply(op ...envDocItemBuilderOp) *envDocItemsBuilder {
	for _, o := range op {
		o(b)
	}
	return b
}

func (b *envDocItemsBuilder) items() []*EnvDocItem {
	items := make([]*EnvDocItem, len(b.names))
	for i, name := range b.names {
		item := &EnvDocItem{
			Doc:      b.doc,
			Opts:     b.opts,
			Children: b.children,
		}
		if name != "" {
			item.Name = fmt.Sprintf("%s%s", b.envPrefix, name)
		}
		items[i] = item
	}
	return items
}

func (b *envDocItemsBuilder) GoString() string {
	return fmt.Sprintf("envDocItemsBuilder{envPrefix: %q, names: %q, doc: %q, opts: %v, children: %v}",
		b.envPrefix, b.names, b.doc, b.opts, b.children)
}

func debugBuilder(l *log.Logger, prefix string, b *envDocItemsBuilder) {
	l.Printf("%s: %s", prefix, b.GoString())
}

func strConcat(s ...string) string {
	var b strings.Builder
	var size int
	for _, v := range s {
		size += len(v)
	}
	b.Grow(size)
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}
