package bottomtabs

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"fmt"
	"html/template"
)

var (
	//go:embed bottomtabs.html
	bottomTabsHTML string
	comp           = component.New("bottomtabs.html", bottomTabsHTML)
)

type Tab struct {
	Icon    string
	Label   string
	Href    string
	Attrs   template.HTMLAttr
	Active  bool
	Primary bool
}

type Options struct {
	Tabs []Tab
}

type templateOptions struct {
	Tabs  []template.HTML
	Style template.HTMLAttr
}

func Render(options Options) (template.HTML, error) {
	return comp.Render(toTemplateOptions(options))
}

func MustRender(options Options) template.HTML {
	return comp.MustRender(toTemplateOptions(options))
}

func toTemplateOptions(options Options) templateOptions {
	tabs := make([]template.HTML, len(options.Tabs))
	for i, t := range options.Tabs {
		tabs[i] = t.render()
	}
	return templateOptions{
		Tabs:  tabs,
		Style: template.HTMLAttr(fmt.Sprintf("style=\"grid-template-columns: repeat(%d, 1fr)\"", len(options.Tabs))),
	}
}

func (t Tab) render() template.HTML {
	tag := "button"
	if t.Href != "" {
		tag = "a"
	}

	iconClass := "icon-wght-200"
	if t.Active {
		iconClass = "icon-fill-1 icon-wght-800 drop-shadow"
	}
	iconHTML := fmt.Sprintf(`<span class="icon %s">%s</span>`, iconClass, t.Icon)

	stateClass := "border-outline-variant"
	labelClass := ""
	if t.Active {
		stateClass = "text-on-primary-container bg-primary-container"
		labelClass = ` class="font-bold"`
	}

	href := ""
	if t.Href != "" {
		href = fmt.Sprintf(` href="%s"`, t.Href)
	}

	attrs := ""
	if t.Attrs != "" {
		attrs = fmt.Sprintf(` %s`, t.Attrs)
	}

	return template.HTML(fmt.Sprintf(`
<%s%s%s class="flex flex-col items-center justify-center border-t-4 gap-0.5 p-2 %s">
  %s
  <label%s>%s</label>
</%s>`, tag, href, attrs, stateClass, iconHTML, labelClass, t.Label, tag))
}
