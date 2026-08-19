package datetimeinput

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/field"
	_ "embed"
	"html/template"
)

// Type maps to a native HTML input type, so the browser's own picker and
// validation apply.
type Type int

const (
	Date Type = iota
	Time
	DateTimeLocal
)

var typeAttr = map[Type]string{
	Date:          "date",
	Time:          "time",
	DateTimeLocal: "datetime-local",
}

var (
	//go:embed datetimeinput.html
	datetimeinputHTML string
	datetimeinputTpl  = component.New("datetimeinput.html", datetimeinputHTML)
)

type Options struct {
	Id    string
	Name  string
	Label string
	Type  Type
	Value string
	// Min, Max, and Step are passed straight through as the input's
	// min/max/step attributes, in whatever format the Type expects (e.g.
	// "2024-01-01" for Date, "09:00" for Time). Left blank, none is set.
	Min      string
	Max      string
	Step     string
	Hint     string
	Error    string
	Required bool
	Attrs    template.HTMLAttr
}

type templateOptions struct {
	Id          string
	Name        string
	Type        string
	Value       string
	Min         string
	Max         string
	Step        string
	Required    bool
	Attrs       template.HTMLAttr
	DescribedBy template.HTMLAttr
	Invalid     template.HTMLAttr
}

func Render(options Options) (template.HTML, error) {
	described := field.Describe(options.Id, options.Hint, options.Error)

	input, err := datetimeinputTpl.Render(templateOptions{
		Id:          options.Id,
		Name:        options.Name,
		Type:        typeAttr[options.Type],
		Value:       options.Value,
		Min:         options.Min,
		Max:         options.Max,
		Step:        options.Step,
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
