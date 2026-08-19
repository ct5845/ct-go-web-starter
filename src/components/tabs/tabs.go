package tabs

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

// Axis is the direction tabs lay out in.
type Axis int

const (
	Horizontal Axis = iota
	Vertical
)

var (
	//go:embed tabs.html
	tabsHTML string
	//go:embed tabs.js
	tabsJS  string
	tabsTpl = component.WithIIFE("tabs.html", tabsHTML, tabsJS)

	//go:embed tab.html
	tabHTML string
	tabTpl  = component.New("tab.html", tabHTML)
)

type Tab struct {
	Label  string
	Href   string
	Attrs  template.HTMLAttr
	Active bool
}

type Options struct {
	Axis Axis
	Tabs []Tab
}

type templateOptions struct {
	Class string
	Tabs  []template.HTML
}

var axisClass = map[Axis]string{
	Horizontal: "flex-row overflow-x-auto rounded-2xl border border-outline-variant",
	Vertical:   "flex-col overflow-y-auto rounded-2xl border border-outline-variant",
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

	return tabsTpl.Render(templateOptions{
		Class: "flex gap-2 bg-container-lowest " + axisClass[options.Axis],
		Tabs:  tabs,
	})
}

type tabProps struct {
	Href        string
	Attrs       template.HTMLAttr
	AriaCurrent template.HTMLAttr
	Class       string
	Label       string
}

func (t Tab) render() (template.HTML, error) {
	class := "btn rounded-none! whitespace-nowrap text-on-surface-variant"
	var ariaCurrent template.HTMLAttr
	if t.Active {
		class = "btn rounded-none! whitespace-nowrap bg-secondary-container text-on-secondary-container font-semibold!"
		ariaCurrent = `aria-current="page"`
	}

	return tabTpl.Render(tabProps{
		Href:        t.Href,
		Attrs:       t.Attrs,
		AriaCurrent: ariaCurrent,
		Class:       class,
		Label:       t.Label,
	})
}
