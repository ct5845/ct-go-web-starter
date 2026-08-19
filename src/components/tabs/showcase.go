package tabs

import (
	"ct-go-web-starter/src/components/demo"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "tabs",
	Title:       "Tabs",
	Source:      "components/tabs",
	Description: "A strip of navigation links, laid out horizontally or vertically.",
	Render:      renderShowcase,
}

func renderShowcase(demo.Request) (template.HTML, error) {
	horizontal, err := Render(Options{
		Axis: Horizontal,
		Tabs: []Tab{
			{Label: "Overview", Href: "#", Active: true},
			{Label: "Details", Href: "#"},
			{Label: "Settings", Href: "#"},
		},
	})
	if err != nil {
		return "", err
	}

	vertical, err := Render(Options{
		Axis: Vertical,
		Tabs: []Tab{
			{Label: "Overview", Href: "#", Active: true},
			{Label: "Details", Href: "#"},
			{Label: "Settings", Href: "#"},
		},
	})
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Note:  "Axis decides flex direction. Both are inline elements sized to their content.",
		Class: "max-w-2xl border border-outline-variant rounded-xl overflow-hidden flex flex-col gap-6 p-4",
		Content: template.HTML(`<div>` + string(horizontal) + `</div>` +
			`<div>` + string(vertical) + `</div>`),
	})
}
