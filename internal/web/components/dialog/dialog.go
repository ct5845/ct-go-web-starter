package dialog

import (
	"ct-go-web-starter/internal/web/components/component"
	_ "embed"
	"html/template"
)

var (
	//go:embed dialog.html
	dialogHTML string
	dialogTpl  = component.New("dialog.html", dialogHTML)
)

type Options struct {
	Id      string
	Content template.HTML
	Label   string
}

func Render(options Options) (template.HTML, error) {
	return dialogTpl.Render(options)
}
