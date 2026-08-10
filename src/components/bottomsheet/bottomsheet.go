package bottomsheet

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

var (
	//go:embed bottomsheet.html
	bottomsheetHTML string
	bottomsheetTpl  = component.New("bottomsheet.html", bottomsheetHTML)
)

type Options struct {
	Id      string
	Content template.HTML
	Label   string
}

func Render(options Options) (template.HTML, error) {
	return bottomsheetTpl.Render(options)
}
