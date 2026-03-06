package home

import (
	_ "embed"
	"io"
	"log/slog"
	"net/http"

	"ct-go-web-starter/src/shared/component"
	"ct-go-web-starter/src/shared/components/footer"
	"ct-go-web-starter/src/shared/components/header"
	"ct-go-web-starter/src/shared/templates"
)

//go:embed home.html
var homeHTML string
var tmpl = component.New("home.html", homeHTML)

func Handler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Rendering index page", "pages", "index", "path", r.URL.Path)
	if r.URL.Path != "/" {
		slog.Warn("Path not found", "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	headerHTML, err := header.Render("Title", "CT Go Web Starter")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	footerHTML, err := footer.Render()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	contentHTML := tmpl.MustRender()

	page, err := templates.Render(
		"Title", "CT Go Web Starter",
		"HeaderHTML", headerHTML,
		"ContentHTML", contentHTML,
		"FooterHTML", footerHTML,
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	slog.Info("Index page rendered successfully", "pages", "index")
	io.WriteString(w, string(page))
}
