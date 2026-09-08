package component

import (
	"bytes"
	"fmt"
	"html/template"
	texttemplate "text/template"
)

type component struct {
	name           string
	template       *template.Template
	scriptTemplate *texttemplate.Template
}

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

// The JS template uses <<< >>> delimiters to avoid conflicts with Go/JavaScript templates.
func withScript(name, htmlTemplate, jsTemplate string) *component {
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

// WithIIFE keeps the JS template's declarations out of the global scope.
func WithIIFE(name, htmlTemplate, jsTemplate string) *component {
	wrappedJS := fmt.Sprintf(`(function() {
%s
})();`, jsTemplate)
	return withScript(name, htmlTemplate, wrappedJS)
}

// WithAlpine defers the JS template until the alpine:init event, so it can
// register Alpine stores and data components.
func WithAlpine(name, htmlTemplate, jsTemplate string) *component {
	alpineWrapper := fmt.Sprintf(`document.addEventListener("alpine:init", () => {
			%s
		});
	`, jsTemplate)
	return withScript(name, htmlTemplate, alpineWrapper)
}

func (c *component) Render(data any) (template.HTML, error) {
	var buf bytes.Buffer
	if err := c.template.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("component %q: %w", c.name, err)
	}
	return template.HTML(buf.String()), nil
}
