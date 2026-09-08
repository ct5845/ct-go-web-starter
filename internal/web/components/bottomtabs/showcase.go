package bottomtabs

import (
	"ct-go-web-starter/internal/web/components/demo"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "bottomtabs",
	Title:       "Bottom tabs",
	Source:      "components/bottomtabs",
	Description: "The mobile tab bar. Normally only visible below the lg breakpoint.",
	Render:      renderShowcase,
}

func renderShowcase(demo.Request) (template.HTML, error) {
	rendered, err := Render(Options{
		Tabs: []Tab{
			{Label: "Home", Href: "#", Icon: "home", Active: true},
			{Label: "Search", Href: "#", Icon: "search"},
			{Label: "More", Icon: "menu"},
		},
	})
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Note:    "Normally pinned to the bottom of the viewport below the lg breakpoint. Shown inline here.",
		Class:   "max-w-md border border-outline-variant rounded-xl overflow-hidden",
		Content: rendered,
	})
}
