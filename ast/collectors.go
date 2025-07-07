package ast

import (
	"sort"
	"strings"

	"github.com/g4s8/envdoc/debug"
)

type RootCollectorOption func(*RootCollector)

func WithFileGlob(glob func(string) bool) RootCollectorOption {
	return func(c *RootCollector) {
		c.fileGlob = glob
	}
}

func WithTypeGlob(glob func(string) bool) RootCollectorOption {
	return func(c *RootCollector) {
		c.typeGlob = glob
	}
}

func WithGoGenDecl(line int, file string) RootCollectorOption {
	return func(c *RootCollector) {
		c.gogenDecl = &struct {
			line int
			file string
		}{
			line: line,
			file: file,
		}
	}
}

var (
	_ interface {
		FileHandler
		TypeHandler
		CommentHandler
	} = (*RootCollector)(nil)
	_ interface {
		DocHandler
		CommentHandler
		FieldHandler
	} = (*TypeCollector)(nil)
	_ FieldHandler = (*FieldCollector)(nil)
)

type RootCollector struct {
	baseDir   string
	fileGlob  func(string) bool
	typeGlob  func(string) bool
	gogenDecl *struct {
		line int
		file string
	}

	// pendingType is true if gogen declaration was specified
	// and the next type will be the expected one
	pendingType bool

	files []*FileSpec
}

var globAcceptAll = func(string) bool { return true }

func NewRootCollector(baseDir string, opts ...RootCollectorOption) *RootCollector {
	c := &RootCollector{
		baseDir:  baseDir,
		fileGlob: globAcceptAll,
		typeGlob: globAcceptAll,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *RootCollector) Files() []*FileSpec {
	// order by file name
	res := make([]*FileSpec, len(c.files))
	copy(res, c.files)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

func (c *RootCollector) onFile(f *FileSpec) interface {
	TypeHandler
	CommentHandler
	ImportHandler
} {
	// convert file name to relative path using baseDir
	// if baseDir is empty or `.` then the file name is used as is.
	name := f.Name
	if c.baseDir != "" && c.baseDir != "." {
		name, _ = strings.CutPrefix(name, c.baseDir)
		name, _ = strings.CutPrefix(name, "/")
		name = "./" + name
	}
	f.Name = name

	if c.fileGlob(f.Name) {
		f.Export = true
	}
	debug.Logf("# COL: file %q, export=%t\n", f.Name, f.Export)
	c.files = append(c.files, f)
	return c
}

func (c *RootCollector) addImport(spec *ImportSpec) {
	currentFile := c.currentFile()
	currentFile.Imports = append(currentFile.Imports, spec)
}

func (c *RootCollector) currentFile() *FileSpec {
	if len(c.files) == 0 {
		panic("emitted type without file")
	}
	return c.files[len(c.files)-1]
}

func (c *RootCollector) onType(tpe *TypeSpec) interface {
	FieldHandler
	CommentHandler
} {
	currentFile := c.currentFile()

	var export bool
	if c.gogenDecl != nil {
		if c.pendingType {
			c.pendingType = false
			export = true
		}
	} else if c.typeGlob(tpe.Name) {
		export = true
	}

	tpe.Export = export
	currentFile.Types = append(currentFile.Types, tpe)
	return &TypeCollector{spec: tpe}
}

func (c *RootCollector) setComment(spec *CommentSpec) {
	currentFile := c.currentFile()

	if c.gogenDecl == nil {
		return
	}
	if c.gogenDecl.file == currentFile.Name && c.gogenDecl.line == spec.Line {
		c.pendingType = true
	}
}

type TypeCollector struct {
	spec *TypeSpec
}

func (c *TypeCollector) setDoc(spec *DocSpec) {
	c.spec.Doc = spec.Doc
}

func (c *TypeCollector) setComment(spec *CommentSpec) {
	if c.spec.Doc != "" {
		c.spec.Doc = spec.Text
	}
}

func (c *TypeCollector) onField(f *FieldSpec) FieldHandler {
	c.spec.Fields = append(c.spec.Fields, f)
	return &FieldCollector{spec: f}
}

type FieldCollector struct {
	spec *FieldSpec
}

func (c *FieldCollector) onField(f *FieldSpec) FieldHandler {
	c.spec.Fields = append(c.spec.Fields, f)
	return &FieldCollector{spec: f}
}
