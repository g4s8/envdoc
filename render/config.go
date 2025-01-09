package render

import "github.com/g4s8/envdoc/types"

type renderItemConfig struct {
	SeparatorFormat  string
	SeparatorDefault string
	OptRequired      string
	OptExpand        string
	OptNonEmpty      string
	OptFromFile      string
	EnvDefaultFormat string
}
type renderConfig struct {
	Item renderItemConfig
	tmpl template
}

var configs = map[types.OutFormat]renderConfig{
	types.OutFormatMarkdown: {
		Item: renderItemConfig{
			SeparatorFormat:  "separated by `%s`",
			SeparatorDefault: "comma-separated",
			OptRequired:      "**required**",
			OptExpand:        "expand",
			OptFromFile:      "from-file",
			OptNonEmpty:      "non-empty",
			EnvDefaultFormat: "default: `%s`",
		},
		tmpl: newTmplText("markdown.tmpl"),
	},
	types.OutFormatHTML: {
		Item: renderItemConfig{
			SeparatorFormat:  `separated by "<code>%s</code>"`,
			SeparatorDefault: "comma-separated",
			OptRequired:      "<strong>required</strong>",
			OptExpand:        "expand",
			OptFromFile:      "from-file",
			OptNonEmpty:      "non-empty",
			EnvDefaultFormat: "default: <code>%s</code>",
		},
		tmpl: newTmplText("html.tmpl"),
	},
	types.OutFormatTxt: {
		Item: renderItemConfig{
			SeparatorFormat:  "separated by `%s`",
			SeparatorDefault: "comma-separated",
			OptRequired:      "required",
			OptExpand:        "expand",
			OptFromFile:      "from-file",
			OptNonEmpty:      "non-empty",
			EnvDefaultFormat: "default: `%s`",
		},
		tmpl: newTmplText("plaintext.tmpl"),
	},
	types.OutFormatEnv: {
		Item: renderItemConfig{
			SeparatorFormat:  "separated by '%s'",
			SeparatorDefault: "comma-separated",
			OptRequired:      "required",
			OptExpand:        "expand",
			OptFromFile:      "from-file",
			OptNonEmpty:      "non-empty",
			EnvDefaultFormat: "default: '%s'",
		},
		tmpl: newTmplText("dotenv.tmpl"),
	},
	types.OutFormatJSON: {
		Item: renderItemConfig{},
		tmpl: newTmplText("json.tmpl"),
	},
}
