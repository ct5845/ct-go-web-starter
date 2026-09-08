package textarea

import (
	"ct-go-web-starter/internal/web/components/demo"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "textarea",
	Title:       "Textarea",
	Group:       "Inputs",
	Source:      "components/textarea",
	Description: "A multi-line text input, wrapped with a label, hint, and error via the field component.",
	Render:      renderShowcase,
}

func renderShowcase(demo.Request) (template.HTML, error) {
	rendered, err := Render(Options{
		Id:          "showcase-textarea",
		Name:        "message",
		Label:       "Message",
		Placeholder: "Say something...",
		Hint:        "Markdown is supported.",
	})
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Class:   "max-w-md",
		Content: rendered,
	})
}
