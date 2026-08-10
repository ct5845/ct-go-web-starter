package nav

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/icon"
	_ "embed"
	"html/template"
)

//go:embed sidebarheader.html
var sidebarheaderHTML string
var sidebarheaderCmp = component.New("sidebarheader.html", sidebarheaderHTML)

func RenderSidebarHeader() (template.HTML, error) {
	return sidebarheaderCmp.Render(struct {
		Logo template.HTML
	}{
		Logo: icon.SVG["favicon"],
	})
}
