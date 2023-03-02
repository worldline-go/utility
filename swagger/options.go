package swagger

type options struct {
	Version     *string
	Host        *string
	BasePath    *string
	Title       *string
	Description *string
	Schemes     []string

	Delims           string
	InfoInstanceName string
	Custom           map[string]interface{}
}

type Option func(*options)

// WithInfoInstanceName sets the name of the swagger instance to be used.
//
// The default name is "swagger".
func WithInfoInstanceName(infoInstanceName string) Option {
	return func(opts *options) {
		opts.InfoInstanceName = infoInstanceName
	}
}

// WithCustom sets the custom data to be used in the swagger template.
//
// Access custom data in the template with {{.Custom.key}}.
func WithCustom(custom map[string]interface{}) Option {
	return func(opts *options) {
		opts.Custom = custom
	}
}

// WithDelims sets the delimiters to be used in the SetInfo function.
//
// The default delimiters are "[[ ]]"
func WithDelims(delims string) Option {
	return func(opts *options) {
		opts.Delims = delims
	}
}

// WithVersion sets the version of the API.
func WithVersion(version string) Option {
	return func(opts *options) {
		opts.Version = &version
	}
}

// WithHost sets the host of the API.
func WithHost(host string) Option {
	return func(opts *options) {
		opts.Host = &host
	}
}

// WithBasePath sets the base path of the API.
func WithBasePath(basePath string) Option {
	return func(opts *options) {
		opts.BasePath = &basePath
	}
}

// WithSchemes sets the schemes of the API.
func WithSchemes(schemes ...string) Option {
	return func(opts *options) {
		opts.Schemes = schemes
	}
}

// WithTitle sets the title of the API.
func WithTitle(title string) Option {
	return func(opts *options) {
		opts.Title = &title
	}
}

// WithDescription sets the description of the API.
func WithDescription(description string) Option {
	return func(opts *options) {
		opts.Description = &description
	}
}
