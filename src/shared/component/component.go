package component

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	texttemplate "text/template"

	_ "ct-go-web-starter/src/infrastructure/config"
)

// component represents a reusable template component
type component struct {
	name           string
	template       *template.Template
	scriptTemplate *texttemplate.Template
}

// New creates a new component with the given name and HTML template string
func New(name, htmlTemplate string) *component {
	slog.Debug("Creating new component", "name", name)
	tmpl, err := template.New(name).Parse(htmlTemplate)
	if err != nil {
		slog.Error("Failed to parse template", "name", name, "error", err)
		panic(err)
	}

	slog.Debug("Component created successfully", "name", name)
	return &component{
		name:     name,
		template: tmpl,
	}
}

// WithJS creates a component with both HTML and JavaScript templates.
// The JS template uses <<< >>> delimiters to avoid conflicts with Go/Alpine templates.
func WithJS(name, htmlTemplate, jsTemplate string) *component {
	slog.Debug("Creating new component with JS", "name", name)

	var scriptTmpl *texttemplate.Template
	if jsTemplate != "" {
		var err error
		scriptTmpl, err = texttemplate.New(name+".js").Delims("<<<", ">>>").Parse(jsTemplate)
		if err != nil {
			slog.Error("Failed to parse JS template", "name", name, "error", err)
			panic(err)
		}
		htmlTemplate += `<script>{{ ComponentJS . }}</script>`
	}

	tmpl, err := template.New(name).
		Funcs(template.FuncMap{
			"ComponentJS": func(data any) template.JS {
				if scriptTmpl == nil {
					return template.JS("")
				}
				var buf bytes.Buffer
				if err := scriptTmpl.Execute(&buf, data); err != nil {
					slog.Error("Failed to execute JS template", "name", name, "error", err)
					return template.JS("")
				}
				return template.JS(buf.String())
			},
		}).
		Parse(htmlTemplate)
	if err != nil {
		slog.Error("Failed to parse template", "name", name, "error", err)
		panic(err)
	}

	slog.Debug("Component with JS created successfully", "name", name)
	return &component{
		name:           name,
		template:       tmpl,
		scriptTemplate: scriptTmpl,
	}
}

// props builds a map[string]any from key/value pairs, matching slog-style variadic args.
// Panics if an odd number of args is passed, or if a key is not a string.
func props(keysAndValues []any) map[string]any {
	if len(keysAndValues) == 0 {
		return nil
	}
	if len(keysAndValues)%2 != 0 {
		panic(fmt.Sprintf("component: Render called with odd number of args (%d)", len(keysAndValues)))
	}
	m := make(map[string]any, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			panic(fmt.Sprintf("component: key at index %d is not a string", i))
		}
		m[key] = keysAndValues[i+1]
	}
	return m
}

// Render executes the component template with the provided key/value pairs and returns the HTML.
// Usage: comp.Render("Title", "Hello", "Count", 42)
func (c *component) Render(keysAndValues ...any) (template.HTML, error) {
	slog.Debug("Executing component template", "name", c.name)
	var buf bytes.Buffer
	err := c.template.Execute(&buf, props(keysAndValues))
	if err != nil {
		slog.Error("Template execution failed", "name", c.name, "error", err)
		return "", err
	}
	slog.Debug("Component template executed successfully", "name", c.name)
	return template.HTML(buf.String()), nil
}

// MustRender executes the component template and panics on error (useful for compile-time safety)
func (c *component) MustRender(keysAndValues ...any) template.HTML {
	html, err := c.Render(keysAndValues...)
	if err != nil {
		slog.Error("Component render failed, panicking", "name", c.name, "error", err)
		panic(err)
	}
	return html
}
