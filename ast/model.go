package ast

import "strings"

//go:generate stringer -type=FieldTypeRefKind -trimprefix=FieldType
type FieldTypeRefKind int

const (
	FieldTypeIdent FieldTypeRefKind = iota
	FieldTypeSelector
	FieldTypePtr
	FieldTypeArray
	FieldTypeMap
	FieldTypeStruct
)

type (
	FileSpec struct {
		Name    string
		Pkg     string
		Imports []*ImportSpec
		Types   []*TypeSpec
		Export  bool // tru if file should be exported
	}

	TypeSpec struct {
		Name   string
		Doc    string
		Fields []*FieldSpec
		Export bool // true if type should be exported
	}

	FieldSpec struct {
		Names   []string
		Doc     string
		Tag     string
		TypeRef FieldTypeRef
		Fields  []*FieldSpec
	}

	FieldTypeRef struct {
		Name string
		Pkg  string
		Kind FieldTypeRefKind
	}

	DocSpec struct {
		Doc string
	}

	CommentSpec struct {
		Text string
		Line int
	}

	ImportSpec struct {
		Name string // aka alias
		Path string
	}
)

type (
	CommentHandler interface {
		setComment(*CommentSpec)
	}

	DocHandler interface {
		setDoc(*DocSpec)
	}

	FieldHandler interface {
		onField(*FieldSpec) FieldHandler
	}

	TypeHandler interface {
		onType(*TypeSpec) interface {
			FieldHandler
			CommentHandler
		}
	}

	FileHandler interface {
		onFile(*FileSpec) interface {
			TypeHandler
			CommentHandler
			ImportHandler
		}
	}

	ImportHandler interface {
		addImport(*ImportSpec)
	}
)

func (tr FieldTypeRef) String() string {
	switch tr.Kind {
	case FieldTypeIdent:
		return tr.Name
	case FieldTypeSelector:
		return tr.Pkg + "." + tr.Name
	case FieldTypePtr:
		return "*" + tr.Name
	case FieldTypeArray:
		return "[]" + tr.Name
	case FieldTypeMap:
		return "map[string]" + tr.Name
	case FieldTypeStruct:
		return "struct"
	}
	return ""
}

func (tr FieldTypeRef) IsBuiltIn() bool {
	switch tr.Name {
	case
		"string",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"float32",
		"float64",
		"bool",
		"byte",
		"rune",
		"complex64",
		"complex128":
		return true

	default:
		return false
	}
}

func (fs *FileSpec) String() string {
	return fs.Name
}

func (ts *TypeSpec) String() string {
	return ts.Name
}

func (fs *FieldSpec) String() string {
	return strings.Join(fs.Names, ", ")
}

func (i *ImportSpec) PathName() string {
	// convert path `github.com/g4s8/envdoc/ast` to `ast`
	if i.Path == "" {
		return i.Path
	}
	parts := strings.Split(i.Path, "/")
	if len(parts) == 0 {
		return i.Path
	}
	return parts[len(parts)-1]
}
