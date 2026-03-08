package footer

import (
	"ct-go-web-starter/src/shared/component"
	_ "embed"
	"html/template"
)

//go:embed footer.html
var footerHTML string
var comp = component.New("footer.html", footerHTML)

func Render(keysAndValues ...any) (template.HTML, error) {
	return comp.Render(keysAndValues...)
}
