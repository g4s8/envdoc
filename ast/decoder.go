package ast

import (
	"strconv"

	"github.com/g4s8/envdoc/debug"
	"github.com/g4s8/envdoc/tags"
	"github.com/g4s8/envdoc/utils"
)

type FieldSpecDecoder struct {
	envPrefix       string
	tagName         string
	tagDefault      string
	tagPrefix       string
	tagSeparator    string
	tagDescription  string
	tagRequired     string
	requiredIfNoDef bool
	useFieldNames   bool
}

func NewFieldSpecDecoder(envPrefix, tagName, tagDefault, tagPrefix, tagSeparator, tagDescription, tagRequired string, requiredIfNoDef, useFieldNames bool) *FieldSpecDecoder {
	return &FieldSpecDecoder{
		envPrefix:       envPrefix,
		tagName:         tagName,
		tagDefault:      tagDefault,
		tagPrefix:       tagPrefix,
		tagSeparator:    tagSeparator,
		tagDescription:  tagDescription,
		tagRequired:     tagRequired,
		requiredIfNoDef: requiredIfNoDef,
		useFieldNames:   useFieldNames,
	}
}

type FieldInfo struct {
	Names       []string
	Required    bool
	Expand      bool
	NonEmpty    bool
	FromFile    bool
	Default     string
	Separator   string
	Description string
}

func (d *FieldSpecDecoder) decodeFieldNames(f *FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	var names []string
	if envName, ok := tag.GetFirst(d.tagName); ok {
		names = []string{envName}
	} else if d.useFieldNames && len(f.Names) > 0 {
		names = make([]string, len(f.Names))
		for i, name := range f.Names {
			names[i] = utils.CamelToSnake(name)
		}
	}
	for i, name := range names {
		names[i] = d.envPrefix + name
	}
	out.Names = names
}

func (d *FieldSpecDecoder) decodeTagValues(_ *FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	if tagValues := tag.GetAll(d.tagName); len(tagValues) > 1 {
		for _, tagValue := range tagValues[1:] {
			switch tagValue {
			case "required":
				out.Required = true
			case "expand":
				out.Expand = true
			case "notEmpty":
				out.Required = true
				out.NonEmpty = true
			case "file":
				out.FromFile = true
			}
		}
	}
}

func (d *FieldSpecDecoder) decodeEnvDefault(_ *FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	if envDefault, ok := tag.GetFirst(d.tagDefault); ok {
		out.Default = envDefault
	} else if d.requiredIfNoDef {
		out.Required = true
	}
}

func (d *FieldSpecDecoder) decodeTagDescription(_ *FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	if d.tagDescription == "" {
		return
	}

	if tagDescription, ok := tag.GetFirst(d.tagDescription); ok {
		out.Description = tagDescription
	}
}

func (d *FieldSpecDecoder) decodeTagRequired(_ *FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	if d.tagRequired == "" {
		return
	}

	if tagRequired, ok := tag.GetFirst(d.tagRequired); ok {
		boolValue, err := strconv.ParseBool(tagRequired)
		if err != nil {
			debug.Logf("# AST: skip required tag[%v] by error: %v\n", tagRequired, err)
			return
		}
		out.Required = boolValue
	}
}

func (d *FieldSpecDecoder) decodeEnvSeparator(f *FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	if envSeparator, ok := tag.GetFirst(d.tagSeparator); ok {
		out.Separator = envSeparator
	} else if f.TypeRef.Kind == FieldTypeArray {
		out.Separator = ","
	}
}

func (d *FieldSpecDecoder) Decode(f *FieldSpec) (res FieldInfo, prefix string) {
	tag := tags.ParseFieldTag(f.Tag)

	d.decodeFieldNames(f, &tag, &res)
	d.decodeTagValues(f, &tag, &res)
	d.decodeEnvDefault(f, &tag, &res)
	d.decodeEnvSeparator(f, &tag, &res)
	d.decodeTagDescription(f, &tag, &res)
	d.decodeTagRequired(f, &tag, &res)

	if envPrefix, ok := tag.GetFirst(d.tagPrefix); ok {
		prefix = d.envPrefix + envPrefix
	}

	return
}
