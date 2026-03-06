package templates

import (
	"ct-go-web-starter/src/shared/component"
	_ "embed"
	"html/template"
	"log/slog"
)

//go:embed page.html
var pageHTML string
var comp = component.New("page.html", pageHTML)

func init() {
	slog.Debug("Page template initialized", "component", "page")
}

func Render(keysAndValues ...any) (template.HTML, error) {
	slog.Debug("Rendering page template", "component", "page")
	result, err := comp.Render(keysAndValues...)
	if err != nil {
		slog.Error("Failed to render page template", "error", err)
		return "", err
	}
	slog.Debug("Page template rendered successfully", "component", "page")
	return result, nil
}
