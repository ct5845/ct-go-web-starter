package sidebar

import (
	"ct-go-web-starter/src/components/demo"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "sidebar",
	Title:       "Sidebar",
	Source:      "components/sidebar",
	Description: "The desktop navigation rail, with grouped items and an active state.",
	Render:      renderShowcase,
}

func renderShowcase(demo.Request) (template.HTML, error) {
	rendered, err := Render(Options{
		Groups: []Group{
			{Items: []Item{
				{Label: "Home", Href: "#", Icon: "home", Active: true},
				{Label: "Search", Href: "#", Icon: "search"},
			}},
			{Title: "Group", Items: []Item{
				{Label: "Showcase", Href: "#", Icon: "widgets"},
				{Label: "Settings", Href: "#", Icon: "widgets"},
			}},
		},
	})
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Note:    "Rendered at its natural width. In a page this sits in the aside of layoutswitch, visible from the lg breakpoint up.",
		Class:   "w-64 border border-outline-variant rounded-xl overflow-hidden",
		Content: rendered,
	})
}
