package header

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

//go:embed header.html
var headerHTML string
var comp = component.New("header.html", headerHTML)

func Render(keysAndValues ...any) (template.HTML, error) {
	return comp.Render(keysAndValues...)
}
