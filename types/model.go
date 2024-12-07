package types

// OutFormat is an output format for the documentation.
type OutFormat string

const (
	OutFormatMarkdown      OutFormat = "markdown"
	OutFormatMarkdownTable OutFormat = "markdown-table"
	OutFormatHTML          OutFormat = "html"
	OutFormatTxt           OutFormat = "plaintext"
	OutFormatEnv           OutFormat = "dotenv"
)

// EnvDocItem is a documentation item for one environment variable.
type EnvDocItem struct {
	// Name of the environment variable.
	Name string
	// Doc is a documentation text for the environment variable.
	Doc string
	// Opts is a set of options for environment variable parsing.
	Opts EnvVarOptions
	// Type is a type of environment variable.
	Type string
	// Children is a list of child environment variables.
	Children []*EnvDocItem
}

type EnvScope struct {
	// Name of the scope.
	Name string
	// Doc is a documentation text for the scope.
	Doc string
	// Vars is a list of environment variables.
	Vars []*EnvDocItem
}

// EnvVarOptions is a set of options for environment variable parsing.
type EnvVarOptions struct {
	// Separator is a separator for array types.
	Separator string
	// Required is a flag that enables required check.
	Required bool
	// Expand is a flag that enables environment variable expansion.
	Expand bool
	// NonEmpty is a flag that enables non-empty check.
	NonEmpty bool
	// FromFile is a flag that enables reading environment variable from a file.
	FromFile bool
	// Default is a default value for the environment variable.
	Default string
	// Description is a description value for the environment variable.
	Description string
}
