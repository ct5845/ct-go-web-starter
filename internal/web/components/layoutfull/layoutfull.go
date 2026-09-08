package layoutfull

import (
	"ct-go-web-starter/internal/web/components/component"
	"ct-go-web-starter/internal/web/components/page"
	"ct-go-web-starter/internal/web/components/pageloader"
	_ "embed"
	"html/template"
)

//go:embed layoutfull.html
var layoutHTML string
var comp = component.New("layoutfull.html", layoutHTML)

type Options struct {
	Content template.HTML
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
