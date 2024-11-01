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

func (d *FieldSpecDecoder) Decode(f *FieldSpec) (res FieldInfo, prefix string) {
	tag := tags.ParseFieldTag(f.Tag)

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

	res.Names = names
	if tagValues := tag.GetAll(d.tagName); len(tagValues) > 1 {
		for _, tagValue := range tagValues[1:] {
			switch tagValue {
			case "required":
				res.Required = true
			case "expand":
				res.Expand = true
			case "notEmpty":
				res.Required = true
				res.NonEmpty = true
			case "file":
				res.FromFile = true
			}
		}
	}

	if envDefault, ok := tag.GetFirst(d.tagDefault); ok {
		res.Default = envDefault
	} else if !ok && d.requiredIfNoDef {
		res.Required = true
	}

	if envSeparator, ok := tag.GetFirst("envSeparator"); ok {
		res.Separator = envSeparator
	} else if f.TypeRef.Kind == FieldTypeArray {
		res.Separator = ","
	}

	if envPrefix, ok := tag.GetFirst("envPrefix"); ok {
		prefix = d.envPrefix + envPrefix
	}

	return
}
