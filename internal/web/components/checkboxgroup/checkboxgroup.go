package checkboxgroup

import (
	"ct-go-web-starter/internal/web/components/component"
	_ "embed"
	"html/template"
)

var (
	//go:embed checkboxgroup.html
	checkboxgroupHTML string
	checkboxgroupTpl  = component.New("checkboxgroup.html", checkboxgroupHTML)

	//go:embed option.html
	optionHTML string
	optionTpl  = component.New("option.html", optionHTML)
)

// Option is one choice in the group. Id must be unique across the page, since
// it backs the hidden checkbox input and its label's for attribute.
type Option struct {
	Id      string
	Value   string
	Label   string
	Checked bool
}

// Options configures a multi-choice group of up to about 14 options,
// rendered as a segmented row of pills in the style of tabs and radiogroup.
// Name is the form field name shared by every option's checkbox input; each
// checked option submits as its own Name/Value pair.
type Options struct {
	Name     string
	Label    string
	Options  []Option
	Hint     string
	Error    string
	Required bool
}

type templateOptions struct {
	GroupId     string
	Label       string
	Required    bool
	Hint        string
	HintId      template.HTMLAttr
	Error       string
	ErrorId     template.HTMLAttr
	DescribedBy template.HTMLAttr
	Invalid     template.HTMLAttr
	Options     []template.HTML
}

func Render(options Options) (template.HTML, error) {
	rendered := make([]template.HTML, len(options.Options))
	for i, o := range options.Options {
		r, err := o.render(options.Name)
		if err != nil {
			return "", err
		}
		rendered[i] = r
	}

	groupId := "checkboxgroup-" + options.Name
	hintId := groupId + "-hint"
	errorId := groupId + "-error"

	var describedBy []byte
	if options.Hint != "" {
		describedBy = append(describedBy, hintId...)
	}
	if options.Error != "" {
		if len(describedBy) > 0 {
			describedBy = append(describedBy, ' ')
		}
		describedBy = append(describedBy, errorId...)
	}

	opts := templateOptions{
		GroupId:  groupId,
		Label:    options.Label,
		Required: options.Required,
		Hint:     options.Hint,
		Error:    options.Error,
		Options:  rendered,
	}
	if options.Hint != "" {
		opts.HintId = template.HTMLAttr(`id="` + hintId + `"`)
	}
	if options.Error != "" {
		opts.ErrorId = template.HTMLAttr(`id="` + errorId + `"`)
		opts.Invalid = `data-invalid="true"`
	}
	if len(describedBy) > 0 {
		opts.DescribedBy = template.HTMLAttr(`aria-describedby="` + string(describedBy) + `"`)
	}

	return checkboxgroupTpl.Render(opts)
}

type optionProps struct {
	Id      string
	Name    string
	Value   string
	Label   string
	Checked bool
}

func (o Option) render(name string) (template.HTML, error) {
	return optionTpl.Render(optionProps{
		Id:      o.Id,
		Name:    name,
		Value:   o.Value,
		Label:   o.Label,
		Checked: o.Checked,
	})
}
