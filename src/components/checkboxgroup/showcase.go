package checkboxgroup

import (
	"ct-go-web-starter/src/components/demo"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "checkboxgroup",
	Title:       "Checkbox group",
	Source:      "components/checkboxgroup",
	Description: "A multi-choice group of native checkbox inputs, styled as a segmented row of pills in the style of tabs and radiogroup.",
	Render:      renderShowcase,
}

func renderShowcase(demo.Request) (template.HTML, error) {
	toppings, err := Render(Options{
		Name:  "toppings",
		Label: "Toppings",
		Hint:  "Pick as many as you like.",
		Options: []Option{
			{Id: "showcase-checkboxgroup-toppings-cheese", Value: "cheese", Label: "Cheese", Checked: true},
			{Id: "showcase-checkboxgroup-toppings-pepperoni", Value: "pepperoni", Label: "Pepperoni", Checked: true},
			{Id: "showcase-checkboxgroup-toppings-mushroom", Value: "mushroom", Label: "Mushroom"},
			{Id: "showcase-checkboxgroup-toppings-olive", Value: "olive", Label: "Olive"},
		},
	})
	if err != nil {
		return "", err
	}

	invalid, err := Render(Options{
		Name:     "interests",
		Label:    "Interests",
		Required: true,
		Error:    "Choose at least one interest.",
		Options: []Option{
			{Id: "showcase-checkboxgroup-interests-music", Value: "music", Label: "Music"},
			{Id: "showcase-checkboxgroup-interests-sports", Value: "sports", Label: "Sports"},
			{Id: "showcase-checkboxgroup-interests-art", Value: "art", Label: "Art"},
			{Id: "showcase-checkboxgroup-interests-tech", Value: "tech", Label: "Tech"},
		},
	})
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Note:    "Options wrap onto multiple lines past a handful of choices; the group scrolls only if it can't wrap. Any number of options can be checked at once. Checking an Interest clears its server error immediately, ahead of any re-validation.",
		Class:   "max-w-md flex flex-col gap-4",
		Content: toppings + invalid,
	})
}
