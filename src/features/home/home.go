package home

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/layouttabbed"
	"ct-go-web-starter/src/components/page"
	"ct-go-web-starter/src/features/nav"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

//go:embed home.html
var homeHTML string
var homeTmpl = component.New("home.html", homeHTML)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", HandleGet)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	page, err := render()

	if err != nil {
		slog.Error("Failed to render home page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(page))
}

func render() (template.HTML, error) {
	content, err := homeTmpl.Render(map[string]any{
		"Title":       "CT Go Web Starter",
		"Description": "A modern Go web application starter with HTMX, Alpine.js, and TailwindCSS",
	})
	if err != nil {
		return "", fmt.Errorf("home page: render content: %w", err)
	}

	bottomTabs, err := nav.Render("home")
	if err != nil {
		return "", fmt.Errorf("home page: render bottom tabs: %w", err)
	}

	return layouttabbed.RenderPage(page.Options{
		Title:           "CT Go Web Starter",
		MetaDescription: "A modern Go web application starter with HTMX, Alpine.js, and TailwindCSS",
	}, layouttabbed.Options{
		Content:    content,
		BottomTabs: bottomTabs,
	})
}
