package bottomsheet

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

var (
	//go:embed bottomsheet.html
	bottomSheetHTML string
	comp            = component.New("bottomsheet.html", bottomSheetHTML)
)

type Options struct {
	Id      string
	Content template.HTML
}

func Render(options Options) (template.HTML, error) {
	return comp.Render(options)
}

func MustRender(options Options) template.HTML {
	return comp.MustRender(options)
}
