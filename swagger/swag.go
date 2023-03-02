package swagger

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/swaggo/swag"
)

var (
	defaultInfoInstanceName = "swagger"
	defaultDelims           = "[[ ]]"
)

type swaggerWrapper struct {
	*swag.Spec
	// Custom for customize swagger with custom function.
	Custom map[string]interface{}
}

func SetInfo(opts ...Option) error {
	options := options{
		InfoInstanceName: defaultInfoInstanceName,
		Delims:           defaultDelims,
	}

	for _, opt := range opts {
		opt(&options)
	}

	delims := strings.Fields(strings.ReplaceAll(options.Delims, ",", " "))
	if len(delims) != 2 { //nolint:gomnd // 2 is the number of delimiters
		return fmt.Errorf("invalid delims: %s", options.Delims)
	}

	spec, ok := swag.GetSwagger(options.InfoInstanceName).(*swag.Spec)
	if !ok {
		return fmt.Errorf("failed to get swagger spec: [%s]", options.InfoInstanceName)
	}

	// modify the swagger template
	tpl, err := template.New("swagger_wrapper").Delims(delims[0], delims[1]).Parse(spec.SwaggerTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse swagger template: %w", err)
	}

	// set the custom
	wrapper := swaggerWrapper{
		Spec:   spec,
		Custom: options.Custom,
	}

	var doc bytes.Buffer
	if err = tpl.Execute(&doc, wrapper); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// set the parameters
	spec.SwaggerTemplate = doc.String()

	if options.Version != nil {
		spec.Version = *options.Version
	}

	if options.Host != nil {
		spec.Host = *options.Host
	}

	if options.BasePath != nil {
		spec.BasePath = *options.BasePath
	}

	if options.Title != nil {
		spec.Title = *options.Title
	}

	if options.Description != nil {
		spec.Description = *options.Description
	}

	if options.Schemes != nil {
		spec.Schemes = options.Schemes
	}

	return nil
}
