package dropdown

import (
	"ct-go-web-starter/internal/web/components/demo"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "dropdown",
	Title:       "Dropdown",
	Group:       "Inputs",
	Source:      "components/dropdown",
	Description: "A native select for a small, fixed set of choices.",
	Render:      renderShowcase,
}

func renderShowcase(demo.Request) (template.HTML, error) {
	country, err := Render(Options{
		Id:          "showcase-dropdown-country",
		Name:        "country",
		Label:       "Country",
		Placeholder: "Choose a country",
		Hint:        "Used to set your default currency.",
		Items: []Item{
			{Value: "uk", Label: "United Kingdom"},
			{Value: "us", Label: "United States"},
			{Value: "de", Label: "Germany"},
			{Value: "fr", Label: "France"},
		},
	})
	if err != nil {
		return "", err
	}

	plan, err := Render(Options{
		Id:       "showcase-dropdown-plan",
		Name:     "plan",
		Label:    "Plan",
		Required: true,
		Error:    "Choose a plan to continue.",
		Items: []Item{
			{Value: "free", Label: "Free"},
			{Value: "pro", Label: "Pro", Selected: true},
			{Value: "team", Label: "Team"},
		},
	})
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Note:    "A native <select>, styled to match the other field controls. Country has a placeholder option; Plan has one pre-selected. Best for a handful of options — for a large option set, a server-driven picker fits better than a single native select.",
		Class:   "max-w-md flex flex-col gap-4",
		Content: country + plan,
	})
}
