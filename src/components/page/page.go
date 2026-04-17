package page

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

//go:embed page.html
var pageHTML string
var comp = component.New("page.html", pageHTML)

func Render(keysAndValues ...any) (template.HTML, error) {
	return comp.Render(keysAndValues...)
}
