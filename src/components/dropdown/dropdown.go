package dropdown

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/field"
	_ "embed"
	"html/template"
)

var (
	//go:embed dropdown.html
	dropdownHTML string
	dropdownTpl  = component.New("dropdown.html", dropdownHTML)
)

type Item struct {
	Value    string
	Label    string
	Selected bool
}

// Options renders a native <select> for a small, fixed set of choices. For
// option sets too large to list up front, use a server-driven picker instead.
type Options struct {
	Id          string
	Name        string
	Label       string
	Placeholder string
	Items       []Item
	Hint        string
	Error       string
	Required    bool
}

type templateOptions struct {
	Id          string
	Name        string
	Placeholder string
	Items       []Item
	Required    bool
	DescribedBy template.HTMLAttr
	Invalid     template.HTMLAttr
}

func Render(options Options) (template.HTML, error) {
	described := field.Describe(options.Id, options.Hint, options.Error)

	input, err := dropdownTpl.Render(templateOptions{
		Id:          options.Id,
		Name:        options.Name,
		Placeholder: options.Placeholder,
		Items:       options.Items,
		Required:    options.Required,
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
