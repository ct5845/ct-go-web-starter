package textarea

import (
	"ct-go-web-starter/internal/web/components/component"
	"ct-go-web-starter/internal/web/components/field"
	_ "embed"
	"html/template"
)

var (
	//go:embed textarea.html
	textareaHTML string
	textareaTpl  = component.New("textarea.html", textareaHTML)
)

type Options struct {
	Id          string
	Name        string
	Label       string
	Value       string
	Placeholder string
	Rows        int
	Hint        string
	Error       string
	Required    bool
	Attrs       template.HTMLAttr
}

type templateOptions struct {
	Id          string
	Name        string
	Value       string
	Placeholder string
	Rows        int
	Required    bool
	Attrs       template.HTMLAttr
	DescribedBy template.HTMLAttr
	Invalid     template.HTMLAttr
}

func Render(options Options) (template.HTML, error) {
	described := field.Describe(options.Id, options.Hint, options.Error)

	rows := options.Rows
	if rows == 0 {
		rows = 4
	}

	input, err := textareaTpl.Render(templateOptions{
		Id:          options.Id,
		Name:        options.Name,
		Value:       options.Value,
		Placeholder: options.Placeholder,
		Rows:        rows,
		Required:    options.Required,
		Attrs:       options.Attrs,
		DescribedBy: described.DescribedBy,
		Invalid:     described.Invalid,
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
