package ast

import (
	"github.com/g4s8/envdoc/tags"
	"github.com/g4s8/envdoc/utils"
)

type FieldSpecDecoder struct {
	envPrefix       string
	tagName         string
	tagDefault      string
	requiredIfNoDef bool
	useFieldNames   bool
}

func NewFieldSpecDecoder(envPrefix string, tagName string, tagDefault string, requiredIfNoDef bool, useFieldNames bool) *FieldSpecDecoder {
	return &FieldSpecDecoder{
		envPrefix:       envPrefix,
		tagName:         tagName,
		tagDefault:      tagDefault,
		requiredIfNoDef: requiredIfNoDef,
		useFieldNames:   useFieldNames,
	}
}

type FieldInfo struct {
	Names     []string
	Required  bool
	Expand    bool
	NonEmpty  bool
	FromFile  bool
	Default   string
	Separator string
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
	} else if !ok && d.requiredIfNoDef {
		out.Required = true
	}
}

func (d *FieldSpecDecoder) decodeEnvSeparator(f *FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	if envSeparator, ok := tag.GetFirst("envSeparator"); ok {
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

	if envPrefix, ok := tag.GetFirst("envPrefix"); ok {
		prefix = d.envPrefix + envPrefix
	}

	return
}
