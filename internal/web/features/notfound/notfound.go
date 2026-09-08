package notfound

import (
	"ct-go-web-starter/internal/infrastructure/reqlog"
	"ct-go-web-starter/internal/web/components/component"
	"ct-go-web-starter/internal/web/components/layoutswitch"
	"ct-go-web-starter/internal/web/components/page"
	"ct-go-web-starter/internal/web/features/nav"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

//go:embed notfound.html
var notfoundHTML string
var notfoundTmpl = component.New("notfound.html", notfoundHTML)

const (
	title       = "Page not found"
	description = "The page you asked for does not exist."
)

// The pattern is the bare "/" catch-all, so it answers anything no other route
// claimed. Home registers "GET /{$}" to keep the root itself.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", HandleGet)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	defer reqlog.Track(r.Context(), "notfound.HandleGet", r.URL.Path)()

	rendered, err := render()
	if err != nil {
		slog.Error("Failed to render not found page", "error", err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, string(rendered))
}

func render() (template.HTML, error) {
	content, err := notfoundTmpl.Render(struct {
		Title       string
		Description string
	}{
		Title:       title,
		Description: description,
	})
	if err != nil {
		return "", fmt.Errorf("not found page: render content: %w", err)
	}

	navigation, err := nav.Render("")
	if err != nil {
		return "", fmt.Errorf("not found page: render navigation: %w", err)
	}

	return layoutswitch.RenderPage(page.Options{
		Title:           title,
		MetaDescription: description,
	}, layoutswitch.Options{
		Content:    content,
		BottomTabs: navigation.Footer,
		SideBar:    navigation.SideBar,
	})
}
