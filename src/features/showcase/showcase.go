package showcase

import (
	"ct-go-web-starter/src/components/bottomsheet"
	"ct-go-web-starter/src/components/bottomtabs"
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/demo"
	"ct-go-web-starter/src/components/icon"
	"ct-go-web-starter/src/components/layoutswitch"
	"ct-go-web-starter/src/components/page"
	"ct-go-web-starter/src/components/pagedlist"
	"ct-go-web-starter/src/components/sidebar"
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
	typographyPage,
	colorsPage,
	spacingPage,
	buttonsPage,
	icon.Showcase,
	menuPage,
	meterPage,
	pagedlist.Showcase,
	sidebar.Showcase,
	bottomtabs.Showcase,
	bottomsheet.Showcase,
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

	//go:embed demolink.html
	demoLinkHTML string
	demoLinkTpl  = component.New("demolink.html", demoLinkHTML)
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

	return renderPage("Showcase", "Every component in this starter, on its own page", content)
}

type demoLinkProps struct {
	Slug        string
	Title       string
	Class       string
	AriaCurrent template.HTMLAttr
}

type demoOptions struct {
	Title       string
	Source      string
	Description string
	Links       []template.HTML
	Content     template.HTML
}

func renderDemoPage(current demo.Page, content template.HTML) (template.HTML, error) {
	links := make([]template.HTML, len(demos))
	for i, d := range demos {
		props := demoLinkProps{
			Slug:  d.Slug,
			Title: d.Title,
			Class: "btn btn-outline whitespace-nowrap text-on-surface-variant",
		}
		if d.Slug == current.Slug {
			props.Class = "btn btn-outline whitespace-nowrap bg-primary-container text-on-primary-container font-bold"
			props.AriaCurrent = `aria-current="page"`
		}

		link, err := demoLinkTpl.Render(props)
		if err != nil {
			return "", fmt.Errorf("showcase: render link %q: %w", d.Slug, err)
		}
		links[i] = link
	}

	body, err := demoTpl.Render(demoOptions{
		Title:       current.Title,
		Source:      current.Source,
		Description: current.Description,
		Links:       links,
		Content:     content,
	})
	if err != nil {
		return "", fmt.Errorf("showcase %q: render demo: %w", current.Slug, err)
	}

	return renderPage(current.Title, current.Description, body)
}

func renderPage(title, description string, content template.HTML) (template.HTML, error) {
	navigation, err := nav.Render("showcase")
	if err != nil {
		return "", fmt.Errorf("showcase page: render navigation: %w", err)
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
