package main

import (
	"github.com/g4s8/envdoc/ast"
	"github.com/g4s8/envdoc/tags"
	"github.com/g4s8/envdoc/types"
	"github.com/g4s8/envdoc/utils"
)

type FieldInfo struct {
	Names     []string
	Required  bool
	Expand    bool
	NonEmpty  bool
	FromFile  bool
	Default   string
	Separator string
}

type FieldDecoder interface {
	Decode(f *ast.FieldSpec) (FieldInfo, string)
}

type FieldDecoderOpts struct {
	EnvPrefix       string
	TagName         string
	TagDefault      string
	RequiredIfNoDef bool
	UseFieldNames   bool
}

func NewFieldDecoder(target types.TargetType, opts FieldDecoderOpts) FieldDecoder {
	switch target {
	case types.TargetTypeCaarlos0:
		return &caarlos0fieldDecoder{opts: opts}
	case types.TargetTypeCleanenv:
		return &cleanenvFieldDecoder{opts: opts}
	default:
		panic("unknown target type")
	}
}

type caarlos0fieldDecoder struct {
	opts FieldDecoderOpts
}

func (d *caarlos0fieldDecoder) decodeFieldNames(f *ast.FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	var names []string
	if envName, ok := tag.GetFirst(d.opts.TagName); ok {
		names = []string{envName}
	} else if d.opts.UseFieldNames && len(f.Names) > 0 {
		names = make([]string, len(f.Names))
		for i, name := range f.Names {
			names[i] = utils.CamelToSnake(name)
		}
	}
	for i, name := range names {
		names[i] = d.opts.EnvPrefix + name
	}
	out.Names = names
}

func (d *caarlos0fieldDecoder) decodeTagValues(_ *ast.FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	if tagValues := tag.GetAll(d.opts.TagName); len(tagValues) > 1 {
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

func (d *caarlos0fieldDecoder) decodeEnvDefault(_ *ast.FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	if envDefault, ok := tag.GetFirst(d.opts.TagDefault); ok {
		out.Default = envDefault
	} else if !ok && d.opts.RequiredIfNoDef {
		out.Required = true
	}
}

func (d *caarlos0fieldDecoder) decodeEnvSeparator(f *ast.FieldSpec, tag *tags.FieldTag, out *FieldInfo) {
	if envSeparator, ok := tag.GetFirst("envSeparator"); ok {
		out.Separator = envSeparator
	} else if f.TypeRef.Kind == ast.FieldTypeArray {
		out.Separator = ","
	}
}

func (d *caarlos0fieldDecoder) Decode(f *ast.FieldSpec) (res FieldInfo, prefix string) {
	tag := tags.ParseFieldTag(f.Tag)

	d.decodeFieldNames(f, &tag, &res)
	d.decodeTagValues(f, &tag, &res)
	d.decodeEnvDefault(f, &tag, &res)
	d.decodeEnvSeparator(f, &tag, &res)

	if envPrefix, ok := tag.GetFirst("envPrefix"); ok {
		prefix = d.opts.EnvPrefix + envPrefix
	}

	return
}

type cleanenvFieldDecoder struct {
	opts FieldDecoderOpts
}

func (d *cleanenvFieldDecoder) Decode(f *ast.FieldSpec) (res FieldInfo, prefix string) {
	tag := tags.ParseFieldTag(f.Tag)

	name, _ := tag.GetFirst("env")

	var required bool
	if envRequired, ok := tag.GetFirst("env-required"); ok {
		required = envRequired == "true"
	}

	var defaultVal string
	if envDefault, ok := tag.GetFirst("env-default"); ok {
		defaultVal = envDefault
	}

	var separator string
	if envSeparator, ok := tag.GetFirst("env-separator"); ok {
		separator = envSeparator
	}

	if envPrefix, ok := tag.GetFirst("env-prefix"); ok {
		prefix = d.opts.EnvPrefix + envPrefix
	}

	res.Names = []string{d.opts.EnvPrefix + name}
	res.Required = required
	res.Default = defaultVal
	res.Separator = separator

	return
}
