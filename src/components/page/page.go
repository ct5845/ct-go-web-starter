package page

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/icon"
	_ "embed"
	"html/template"
)

//go:embed page.html
var pageHTML string
var comp = component.New("page.html", pageHTML)

type Options struct {
	Title           string
	MetaDescription string
	Robots          string
	IconsHref       string
	CanonicalURL    string
	OGImageURL      string
	FaviconURL      string
	Body            template.HTML
}

func Render(options Options) (template.HTML, error) {
	options.IconsHref = icon.IconFontHref

	return comp.Render(options)
}
