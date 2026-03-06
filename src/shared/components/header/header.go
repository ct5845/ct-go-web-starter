package header

import (
	"ct-go-web-starter/src/shared/component"
	_ "embed"
	"html/template"
	"log/slog"
)

//go:embed header.html
var headerHTML string
var comp = component.New("header.html", headerHTML)

func init() {
	slog.Debug("Header component initialized", "component", "header")
}

func Render(keysAndValues ...any) (template.HTML, error) {
	slog.Debug("Rendering header component", "component", "header")
	result, err := comp.Render(keysAndValues...)
	if err != nil {
		slog.Error("Failed to render header component", "error", err)
		return "", err
	}
	slog.Debug("Header component rendered successfully", "component", "header")
	return result, nil
}
