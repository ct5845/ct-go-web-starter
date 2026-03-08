package component

import (
	"bytes"
	"fmt"
	"html/template"
	texttemplate "text/template"

	_ "ct-go-web-starter/src/infrastructure/config"
)

type component struct {
	name           string
	template       *template.Template
	scriptTemplate *texttemplate.Template
}

// New creates a new component with the given name and HTML template string
func New(name, htmlTemplate string) *component {
	tmpl, err := template.New(name).Parse(htmlTemplate)
	if err != nil {
		panic(fmt.Sprintf("component: failed to parse template %q: %v", name, err))
	}
	return &component{
		name:     name,
		template: tmpl,
	}
}

// WithJS creates a component with both HTML and JavaScript templates.
// The JS template uses <<< >>> delimiters to avoid conflicts with Go/Alpine templates.
func WithJS(name, htmlTemplate, jsTemplate string) *component {
	var scriptTmpl *texttemplate.Template
	if jsTemplate != "" {
		var err error
		scriptTmpl, err = texttemplate.New(name+".js").Delims("<<<", ">>>").Parse(jsTemplate)
		if err != nil {
			panic(fmt.Sprintf("component: failed to parse JS template %q: %v", name, err))
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
					return template.JS("")
				}
				return template.JS(buf.String())
			},
		}).
		Parse(htmlTemplate)
	if err != nil {
		panic(fmt.Sprintf("component: failed to parse template %q: %v", name, err))
	}

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
	var buf bytes.Buffer
	if err := c.template.Execute(&buf, props(keysAndValues)); err != nil {
		return "", fmt.Errorf("component %q: %w", c.name, err)
	}
	return template.HTML(buf.String()), nil
}

// MustRender executes the component template and panics on error (useful for compile-time safety)
func (c *component) MustRender(keysAndValues ...any) template.HTML {
	html, err := c.Render(keysAndValues...)
	if err != nil {
		panic(err)
	}
	return html
}
