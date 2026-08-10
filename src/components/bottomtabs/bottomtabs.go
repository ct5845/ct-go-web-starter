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
	//go:embed bottomtabs.js
	bottomTabsJS  string
	bottomTabsTpl = component.WithAlpine("bottomtabs.html", bottomTabsHTML, bottomTabsJS)

	//go:embed tab.html
	tabHTML string
	tabTpl  = component.New("tab.html", tabHTML)
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
	tabs := make([]template.HTML, len(options.Tabs))
	for i, t := range options.Tabs {
		rendered, err := t.render()
		if err != nil {
			return "", err
		}
		tabs[i] = rendered
	}
	return bottomTabsTpl.Render(templateOptions{
		Tabs:  tabs,
		Style: template.HTMLAttr(fmt.Sprintf("style=\"grid-template-columns: repeat(%d, 1fr)\"", len(options.Tabs))),
	})
}

type tabProps struct {
	Href       string
	Attrs      template.HTMLAttr
	Class      string
	Icon       string
	IconClass  string
	Label      string
	LabelClass string
}

func (t Tab) render() (template.HTML, error) {
	iconClass := "icon-wght-200 icon-opsz-20"
	stateClass := "border-outline-variant"
	labelClass := "text-sm"
	if t.Active {
		iconClass = "icon-fill-1 icon-wght-800 icon-opsz-20 drop-shadow"
		stateClass = "text-on-primary-container bg-primary-container"
		labelClass = "font-bold text-sm"
	}

	return tabTpl.Render(tabProps{
		Href:       t.Href,
		Attrs:      t.Attrs,
		Class:      "flex flex-col items-center justify-center border-t-2 p-2 " + stateClass,
		Icon:       t.Icon,
		IconClass:  iconClass,
		Label:      t.Label,
		LabelClass: labelClass,
	})
}
