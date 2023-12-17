package main

import (
	"fmt"
	"strings"
)

var (
	fmtPlain itemFormat = plaintextFormat(0)
	fmtMD    itemFormat = markdownFormat(0)
	fmtHTML  itemFormat = htmlFormat(0)
)

type markdownFormat int

func (markdownFormat) formatItem(item docItem) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("- `%s` ", item.envName))
	var opts []string
	if item.flags&docItemFlagRequired != 0 {
		opts = append(opts, "**required**")
	}
	if item.flags&docItemFlagExpand != 0 {
		opts = append(opts, "expand")
	}
	if item.flags&docItemFlagNonEmpty != 0 {
		opts = append(opts, "not-empty")
	}
	if item.flags&docItemFlagFromFile != 0 {
		opts = append(opts, "from file")
	}
	if item.envDefault != "" {
		opts = append(opts, fmt.Sprintf("default: `%s`", item.envDefault))
	}
	if len(opts) > 0 {
		sb.WriteString(fmt.Sprintf("(%s) ", strings.Join(opts, ", ")))
	}
	sb.WriteString(fmt.Sprintf("- %s\n", item.doc))
	return sb.String()
}

func (markdownFormat) beforeItems() string {
	return "# Environment Variables\n\n"
}

func (markdownFormat) afterItems() string {
	return ""
}

type plaintextFormat int

func (plaintextFormat) formatItem(item docItem) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(" * `%s` ", item.envName))
	var opts []string
	if item.flags&docItemFlagRequired != 0 {
		opts = append(opts, "required")
	}
	if item.flags&docItemFlagExpand != 0 {
		opts = append(opts, "expand")
	}
	if item.flags&docItemFlagNonEmpty != 0 {
		opts = append(opts, "not-empty")
	}
	if item.flags&docItemFlagFromFile != 0 {
		opts = append(opts, "from file")
	}
	if item.envDefault != "" {
		opts = append(opts, fmt.Sprintf("default: `%s`", item.envDefault))
	}
	if len(opts) > 0 {
		sb.WriteString(fmt.Sprintf("(%s) ", strings.Join(opts, ", ")))
	}
	sb.WriteString(fmt.Sprintf("- %s\n", item.doc))
	return sb.String()
}

func (plaintextFormat) beforeItems() string {
	return "ENVIRONMENT VARIABLES\n\n"
}

func (plaintextFormat) afterItems() string {
	return ""
}

type htmlFormat int

func (htmlFormat) formatItem(item docItem) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<li><code>%s</code> ", item.envName))
	var opts []string
	if item.flags&docItemFlagRequired != 0 {
		opts = append(opts, "<strong>required</strong>")
	}
	if item.flags&docItemFlagExpand != 0 {
		opts = append(opts, "expand")
	}
	if item.flags&docItemFlagNonEmpty != 0 {
		opts = append(opts, "not-empty")
	}
	if item.flags&docItemFlagFromFile != 0 {
		opts = append(opts, "from file")
	}
	if item.envDefault != "" {
		opts = append(opts, fmt.Sprintf("default: <code>%s</code>", item.envDefault))
	}
	if len(opts) > 0 {
		sb.WriteString(fmt.Sprintf("(%s) ", strings.Join(opts, ", ")))
	}
	sb.WriteString(fmt.Sprintf("- %s</li>\n", item.doc))
	return sb.String()
}

func (htmlFormat) beforeItems() string {
	return `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Environment Variables</title>
		<style>
		body {
			font-family: sans-serif;
		}
		</style>
	</head>
	<body>
		<h1>Environment Variables</h1>
		<ul>
`
}

func (htmlFormat) afterItems() string {
	return `</ul>
	</body>
</html>
`
}
