package menu

import (
	"ct-go-web-starter/src/components/demo"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "menu",
	Title:       "Popover menu",
	Source:      "components/menu",
	Description: "A menu anchored to its trigger with the native popover attribute and CSS anchor positioning, nesting submenus the same way.",
	Render:      renderShowcase,
}

const showcaseMenuID = "showcase-popover-menu"

func renderShowcase(demo.Request) (template.HTML, error) {
	menu, err := Render(Options{
		Id: showcaseMenuID,
		Items: []Item{
			{Label: "Profile", Icon: "widgets", Href: "#"},
			{Label: "Settings", Icon: "widgets", Href: "#"},
			{
				Label: "Share",
				Icon:  "widgets",
				Items: []Item{
					{Label: "Copy link", Href: "#"},
					{Label: "Email", Href: "#"},
					{
						Label: "Export as",
						Items: []Item{
							{Label: "PDF", Href: "#"},
							{Label: "CSV", Href: "#"},
						},
					},
				},
			},
			{Label: "Sign out", Href: "#"},
		},
	})
	if err != nil {
		return "", err
	}

	trigger, err := demo.Frame(demo.FrameOptions{
		Note:    "Uses the native popover attribute, popovertarget, and CSS anchor positioning (anchor-name/position-anchor/position-area) — no JavaScript. Anchor positioning needs a Chromium browser; elsewhere the popover falls back to the browser's default centering.",
		Content: template.HTML(`<button type="button" class="btn btn-primary" popovertarget="` + showcaseMenuID + `" style="anchor-name: ` + AnchorName(showcaseMenuID) + `">Open menu</button>`),
	})
	if err != nil {
		return "", err
	}

	return trigger + menu, nil
}
