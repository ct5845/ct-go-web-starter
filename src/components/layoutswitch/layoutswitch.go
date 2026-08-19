package layoutswitch

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/page"
	"ct-go-web-starter/src/components/pageloader"
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

type templateOptions struct {
	Options
	PageLoader template.HTML
}

func Render(options Options) (template.HTML, error) {
	pageLoader, err := pageloader.Render()
	if err != nil {
		return "", err
	}

	return comp.Render(templateOptions{Options: options, PageLoader: pageLoader})
}

func RenderPage(pageOptions page.Options, options Options) (template.HTML, error) {
	body, err := Render(options)
	if err != nil {
		return "", err
	}

	pageOptions.Body = body
	return page.Render(pageOptions)
}
