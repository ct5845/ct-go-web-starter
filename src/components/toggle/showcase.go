package toggle

import (
	"ct-go-web-starter/src/components/demo"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "toggle",
	Title:       "Toggle",
	Group:       "Inputs",
	Source:      "components/toggle",
	Description: "A boolean switch, backed by a native checkbox input styled as a toggle.",
	Render:      renderShowcase,
}

func renderShowcase(demo.Request) (template.HTML, error) {
	notifications, err := Render(Options{
		Id:      "showcase-toggle-notifications",
		Name:    "notifications",
		Label:   "Email notifications",
		Checked: true,
		Hint:    "You can change this at any time.",
	})
	if err != nil {
		return "", err
	}

	terms, err := Render(Options{
		Id:       "showcase-toggle-terms",
		Name:     "terms",
		Label:    "I agree to the terms of service",
		Required: true,
		Error:    "You must agree to continue.",
	})
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Note:    "A native checkbox styled as a switch, with the label sitting beside it rather than above. Toggling Terms clears its server error immediately, ahead of any re-validation.",
		Class:   "max-w-md flex flex-col gap-4",
		Content: notifications + terms,
	})
}
