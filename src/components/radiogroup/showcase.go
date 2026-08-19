package radiogroup

import (
	"ct-go-web-starter/src/components/demo"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "radiogroup",
	Title:       "Radio group",
	Group:       "Inputs",
	Source:      "components/radiogroup",
	Description: "A single-choice group of native radio inputs, styled as a segmented row of pills in the style of tabs.",
	Render:      renderShowcase,
}

func renderShowcase(demo.Request) (template.HTML, error) {
	plan, err := Render(Options{
		Name:  "plan",
		Label: "Plan",
		Hint:  "You can change this later.",
		Options: []Option{
			{Id: "showcase-radiogroup-plan-free", Value: "free", Label: "Free", Checked: true},
			{Id: "showcase-radiogroup-plan-pro", Value: "pro", Label: "Pro"},
			{Id: "showcase-radiogroup-plan-team", Value: "team", Label: "Team"},
		},
	})
	if err != nil {
		return "", err
	}

	invalid, err := Render(Options{
		Name:     "shift",
		Label:    "Shift",
		Required: true,
		Error:    "Choose a shift to continue.",
		Options: []Option{
			{Id: "showcase-radiogroup-shift-morning", Value: "morning", Label: "Morning"},
			{Id: "showcase-radiogroup-shift-afternoon", Value: "afternoon", Label: "Afternoon"},
			{Id: "showcase-radiogroup-shift-evening", Value: "evening", Label: "Evening"},
			{Id: "showcase-radiogroup-shift-night", Value: "night", Label: "Night"},
		},
	})
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Note:    "Picking a Shift clears its server error immediately, ahead of any re-validation.",
		Class:   "max-w-md flex flex-col gap-4",
		Content: plan + invalid,
	})
}
