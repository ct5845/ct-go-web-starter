package textinput

import (
	"ct-go-web-starter/internal/web/components/component"
	"ct-go-web-starter/internal/web/components/field"
	_ "embed"
	"html/template"
)

// Type maps to a native HTML input type, so the browser's own keyboard,
// validation, and spinner behaviour apply.
type Type int

const (
	Text Type = iota
	Number
	Password
	Email
)

var typeAttr = map[Type]string{
	Text:     "text",
	Number:   "number",
	Password: "password",
	Email:    "email",
}

var (
	//go:embed textinput.html
	textinputHTML string
	textinputTpl  = component.New("textinput.html", textinputHTML)
)

type Options struct {
	Id          string
	Name        string
	Label       string
	Type        Type
	Value       string
	Placeholder string
	// Autocomplete sets the input's autocomplete attribute. Password fields
	// should pass "current-password" or "new-password" so browsers offer the
	// right saved-credential and generation behaviour; left blank, the field
	// falls back to "off" for every other type.
	Autocomplete string
	Hint         string
	Error        string
	Required     bool
	Attrs        template.HTMLAttr
}

type templateOptions struct {
	Id           string
	Name         string
	Type         string
	Value        string
	Placeholder  string
	Autocomplete string
	Required     bool
	Attrs        template.HTMLAttr
	DescribedBy  template.HTMLAttr
	Invalid      template.HTMLAttr
	Reveal       bool
}

func Render(options Options) (template.HTML, error) {
	described := field.Describe(options.Id, options.Hint, options.Error)

	autocomplete := options.Autocomplete
	if autocomplete == "" {
		autocomplete = "off"
	}

	input, err := textinputTpl.Render(templateOptions{
		Id:           options.Id,
		Name:         options.Name,
		Type:         typeAttr[options.Type],
		Value:        options.Value,
		Placeholder:  options.Placeholder,
		Autocomplete: autocomplete,
		Required:     options.Required,
		Attrs:        options.Attrs,
		DescribedBy:  described.DescribedBy,
		Invalid:      described.Invalid,
		Reveal:       options.Type == Password,
	})
	if err != nil {
		return "", err
	}

	return field.Render(field.Options{
		Id:       options.Id,
		Label:    options.Label,
		Hint:     options.Hint,
		Error:    options.Error,
		Required: options.Required,
		Input:    input,
	})
}
