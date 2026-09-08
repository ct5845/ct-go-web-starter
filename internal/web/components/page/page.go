package page

import (
	"ct-go-web-starter/internal/infrastructure/config"
	"ct-go-web-starter/internal/web/components/component"
	"ct-go-web-starter/internal/web/components/icon"
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
	IsDev           bool
}

func Render(options Options) (template.HTML, error) {
	options.IconsHref = icon.IconFontHref
	options.IsDev = config.AppEnv == "dev"

	return comp.Render(options)
}
