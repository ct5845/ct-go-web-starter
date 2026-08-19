// Package pageloader renders a thin progress bar pinned to the top of the
// viewport, shown while a full page navigation or htmx-boosted request is in
// flight. It has no configuration and is meant to be rendered once, in the
// page chrome, rather than per feature.
package pageloader

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

var (
	//go:embed pageloader.html
	pageloaderHTML string
	//go:embed pageloader.js
	pageloaderJS  string
	pageloaderTpl = component.WithIIFE("pageloader.html", pageloaderHTML, pageloaderJS)
)

func Render() (template.HTML, error) {
	return pageloaderTpl.Render(nil)
}
