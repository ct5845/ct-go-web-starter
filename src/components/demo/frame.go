package demo

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

var (
	//go:embed frame.html
	frameHTML string
	frameTpl  = component.New("frame.html", frameHTML)
)

// FrameOptions wraps a demo in the showcase's shared chrome. Note explains
// anything the isolated rendering cannot show — most often where the component
// actually sits in a real page. Class constrains the demo to a realistic size.
type FrameOptions struct {
	Note    string
	Class   string
	Content template.HTML
}

func Frame(options FrameOptions) (template.HTML, error) {
	return frameTpl.Render(options)
}
