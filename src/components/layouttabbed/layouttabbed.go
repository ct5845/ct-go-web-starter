package layouttabbed

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/page"
	_ "embed"
	"html/template"
)

//go:embed layouttabbed.html
var layoutHTML string
var comp = component.New("layouttabbed.html", layoutHTML)

type Tab struct {
}

type Options struct {
	Header     template.HTML
	Content    template.HTML
	BottomTabs template.HTML
}

func Render(options Options) (template.HTML, error) {
	return comp.Render(options)
}

func RenderPage(pageOptions page.Options, options Options) (template.HTML, error) {
	if body, err := Render(options); err != nil {
		return "", err
	} else {
		pageOptions.Body = body
	}

	return page.Render(pageOptions)
}
