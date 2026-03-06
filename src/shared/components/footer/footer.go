package footer

import (
	"ct-go-web-starter/src/shared/component"
	_ "embed"
	"html/template"
	"log/slog"
)

//go:embed footer.html
var footerHTML string
var comp = component.New("footer.html", footerHTML)

func init() {
	slog.Debug("Footer component initialized", "component", "footer")
}

func Render(keysAndValues ...any) (template.HTML, error) {
	slog.Debug("Rendering footer component", "component", "footer")
	result, err := comp.Render(keysAndValues...)
	if err != nil {
		slog.Error("Failed to render footer component", "error", err)
		return "", err
	}
	slog.Debug("Footer component rendered successfully", "component", "footer")
	return result, nil
}
