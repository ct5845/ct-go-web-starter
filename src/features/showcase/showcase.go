package showcase

import (
	"ct-go-web-starter/src/components/bottomsheet"
	"ct-go-web-starter/src/components/bottomtabs"
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/demo"
	"ct-go-web-starter/src/components/dialog"
	"ct-go-web-starter/src/components/icon"
	"ct-go-web-starter/src/components/layoutswitch"
	"ct-go-web-starter/src/components/menu"
	"ct-go-web-starter/src/components/page"
	"ct-go-web-starter/src/components/pagedlist"
	"ct-go-web-starter/src/components/sidebar"
	"ct-go-web-starter/src/components/tabs"
	"ct-go-web-starter/src/features/nav"
	"ct-go-web-starter/src/infrastructure/reqlog"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

// demos is the running order of the showcase. Component demos are declared in
// the component's own package; the stylesheet demos live in styles.go because
// they have no package to sit alongside.
var demos = []demo.Page{
	bottomsheet.Showcase,
	bottomtabs.Showcase,
	buttonsPage,
	colorsPage,
	dialog.Showcase,
	icon.Showcase,
	menu.Showcase,
	menuListPage,
	meterPage,
	pagedlist.Showcase,
	sidebar.Showcase,
	spacingPage,
	tabs.Showcase,
	typographyPage,
}

func lookup(slug string) (demo.Page, bool) {
	for _, d := range demos {
		if d.Slug == slug {
			return d, true
		}
	}
	return demo.Page{}, false
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /showcase", handleIndex)
	mux.HandleFunc("GET /showcase/{slug}", handleDemo)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	defer reqlog.Track(r.Context(), "showcase.handleIndex", "")()

	body, err := renderIndex()
	if err != nil {
		slog.Error("Failed to render showcase index", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(body))
}

func handleDemo(w http.ResponseWriter, r *http.Request) {
	defer reqlog.Track(r.Context(), "showcase.handleDemo", "")()

	current, ok := lookup(r.PathValue("slug"))
	if !ok {
		http.NotFound(w, r)
		return
	}

	content, err := current.Render(demo.Request{
		BaseHref: r.URL.Path,
		Query:    r.URL.Query(),
	})
	if err != nil {
		slog.Error("Failed to render showcase demo", "slug", current.Slug, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// A demo that swaps part of itself via htmx only needs its own markup
	// back; the surrounding page chrome is already on screen.
	if r.Header.Get("HX-Request") == "true" {
		io.WriteString(w, string(content))
		return
	}

	body, err := renderDemoPage(current, content)
	if err != nil {
		slog.Error("Failed to render showcase page", "slug", current.Slug, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(body))
}

var (
	//go:embed index.html
	indexHTML string
	indexTpl  = component.New("index.html", indexHTML)

	//go:embed indexitem.html
	indexItemHTML string
	indexItemTpl  = component.New("indexitem.html", indexItemHTML)

	//go:embed demo.html
	demoHTML string
	demoTpl  = component.New("demo.html", demoHTML)
)

func renderIndex() (template.HTML, error) {
	items := make([]template.HTML, len(demos))
	for i, d := range demos {
		item, err := indexItemTpl.Render(d)
		if err != nil {
			return "", fmt.Errorf("showcase index: render item %q: %w", d.Slug, err)
		}
		items[i] = item
	}

	content, err := indexTpl.Render(struct{ Items []template.HTML }{items})
	if err != nil {
		return "", fmt.Errorf("showcase index: render content: %w", err)
	}

	return renderPage("Showcase", "Every component in this starter, on its own page", content, true)
}

type demoOptions struct {
	Title       string
	Source      string
	Description string
	Nav         template.HTML
	Content     template.HTML
}

func renderDemoPage(current demo.Page, content template.HTML) (template.HTML, error) {
	demoTabs := make([]tabs.Tab, 0, len(demos)+1)
	demoTabs = append(demoTabs, tabs.Tab{Label: "All", Href: "/showcase"})
	for _, d := range demos {
		demoTabs = append(demoTabs, tabs.Tab{
			Label:  d.Title,
			Href:   "/showcase/" + d.Slug,
			Active: d.Slug == current.Slug,
		})
	}

	nav, err := tabs.Render(tabs.Options{Axis: tabs.Horizontal, Tabs: demoTabs})
	if err != nil {
		return "", fmt.Errorf("showcase: render nav: %w", err)
	}

	body, err := demoTpl.Render(demoOptions{
		Title:       current.Title,
		Source:      current.Source,
		Description: current.Description,
		Nav:         nav,
		Content:     content,
	})
	if err != nil {
		return "", fmt.Errorf("showcase %q: render demo: %w", current.Slug, err)
	}

	return renderPage(current.Title, current.Description, body, false)
}

func renderPage(title, description string, content template.HTML, switchlayout bool) (template.HTML, error) {
	navigation, err := nav.Render("showcase")
	if err != nil {
		return "", fmt.Errorf("showcase page: render navigation: %w", err)
	}

	pageOptions := page.Options{
		Title:           title,
		MetaDescription: description,
	}

	var bottomTabs template.HTML

	if switchlayout {
		bottomTabs = navigation.Footer
	} else {
		bottomTabs = ""
	}

	return layoutswitch.RenderPage(pageOptions, layoutswitch.Options{
		Content:    content,
		BottomTabs: bottomTabs,
		SideBar:    navigation.SideBar,
	})
}
