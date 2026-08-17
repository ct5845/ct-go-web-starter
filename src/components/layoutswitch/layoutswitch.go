package layoutswitch

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/page"
	_ "embed"
	"html/template"
)

//go:embed layoutswitch.html
var layoutHTML string
var comp = component.New("layoutswitch.html", layoutHTML)

type Options struct {
	Content    template.HTML
	BottomTabs template.HTML
	SideBar    template.HTML
}

func Render(options Options) (template.HTML, error) {
	return comp.Render(options)
}

func RenderPage(pageOptions page.Options, options Options) (template.HTML, error) {
	body, err := Render(options)
	if err != nil {
		return "", err
	}

	pageOptions.Body = body
	return page.Render(pageOptions)
}
