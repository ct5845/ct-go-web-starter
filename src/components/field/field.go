// Package field wraps a single form input with its label, hint, and error
// text. The wrapper carries Alpine's inputControl(id), registered globally in
// page.html, which dispatches a page-wide "field-change" event on change —
// see page.html for the cross-field convention this enables.
package field

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

var (
	//go:embed field.html
	fieldHTML string
	fieldTpl  = component.New("field.html", fieldHTML)
)

type Options struct {
	Id    string
	Label string
	// NoLabel skips the field's own <label>, for inputs whose template
	// renders its own inline label (e.g. toggle, where the label sits beside
	// the switch rather than above it).
	NoLabel  bool
	Hint     string
	Error    string
	Required bool
	Input    template.HTML
}

type InputAttrs struct {
	DescribedBy template.HTMLAttr
	Invalid     template.HTMLAttr
}

// Describe must be called before rendering the input, using the same id
// passed to Options, so the rendered input and the label/hint/error text
// this package renders around it reference each other correctly.
func Describe(id, hint, errorText string) InputAttrs {
	hintId, errorId := ids(id)

	var describedBy []byte
	if hint != "" {
		describedBy = append(describedBy, hintId...)
	}
	if errorText != "" {
		if len(describedBy) > 0 {
			describedBy = append(describedBy, ' ')
		}
		describedBy = append(describedBy, errorId...)
	}

	attrs := InputAttrs{}
	if len(describedBy) > 0 {
		attrs.DescribedBy = template.HTMLAttr(`aria-describedby="` + string(describedBy) + `"`)
	}
	if errorText != "" {
		attrs.Invalid = `aria-invalid="true"`
	}
	return attrs
}

func ids(id string) (hintId, errorId string) {
	return id + "-hint", id + "-error"
}

type templateOptions struct {
	Id       string
	Label    string
	NoLabel  bool
	Required bool
	Hint     string
	HintId   template.HTMLAttr
	Error    string
	ErrorId  template.HTMLAttr
	Invalid  template.HTMLAttr
	Input    template.HTML
}

func Render(options Options) (template.HTML, error) {
	hintId, errorId := ids(options.Id)

	opts := templateOptions{
		Id:       options.Id,
		Label:    options.Label,
		NoLabel:  options.NoLabel,
		Required: options.Required,
		Hint:     options.Hint,
		Error:    options.Error,
		Input:    options.Input,
	}
	if options.Hint != "" {
		opts.HintId = template.HTMLAttr(`id="` + hintId + `"`)
	}
	if options.Error != "" {
		opts.ErrorId = template.HTMLAttr(`id="` + errorId + `"`)
		opts.Invalid = `data-invalid="true"`
	}

	return fieldTpl.Render(opts)
}
