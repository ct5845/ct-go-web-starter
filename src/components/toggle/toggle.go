package toggle

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/field"
	_ "embed"
	"html/template"
)

var (
	//go:embed toggle.html
	toggleHTML string
	toggleTpl  = component.New("toggle.html", toggleHTML)
)

type Options struct {
	Id       string
	Name     string
	Label    string
	Checked  bool
	Hint     string
	Error    string
	Required bool
}

type templateOptions struct {
	Id          string
	Name        string
	Label       string
	Checked     bool
	Required    bool
	DescribedBy template.HTMLAttr
	Invalid     template.HTMLAttr
}

func Render(options Options) (template.HTML, error) {
	described := field.Describe(options.Id, options.Hint, options.Error)

	input, err := toggleTpl.Render(templateOptions{
		Id:          options.Id,
		Name:        options.Name,
		Label:       options.Label,
		Checked:     options.Checked,
		Required:    options.Required,
		DescribedBy: described.DescribedBy,
		Invalid:     described.Invalid,
	})
	if err != nil {
		return "", err
	}

	return field.Render(field.Options{
		Id:       options.Id,
		Hint:     options.Hint,
		Error:    options.Error,
		Required: options.Required,
		Input:    input,
		NoLabel:  true,
	})
}
