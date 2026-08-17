package showcase

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/layoutswitch"
	"ct-go-web-starter/src/components/page"
	"ct-go-web-starter/src/features/nav"
	"ct-go-web-starter/src/infrastructure/reqlog"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

// demo is one showcase page. Slug matches the component's package name under
// src/components, or the stylesheet name under src/static/styles, so the page
// you are looking at tells you which file to edit.
type demo struct {
	Slug        string
	Title       string
	Source      string
	Description string
	Render      func(r *http.Request) (template.HTML, error)
}

var demos = []demo{
	{"typography", "Typography", "styles/typography.css", "The fluid type scale and base element styles. Resize the window — every size interpolates between its minimum and maximum viewport value.", renderTypography},
	{"colors", "Colours", "styles/theme.css", "Every theme token as a background/foreground pair. Each swatch shows the class names to use.", renderColors},
	{"spacing", "Spacing", "styles/spacing.css", "Clamped spacing utilities that breathe with the viewport, alongside their fixed equivalents.", renderSpacing},
	{"buttons", "Buttons", "styles/button.css", "Button variants, icon buttons, and disabled states.", renderButtons},
	{"icons", "Icons", "components/icon", "The subsetted icon font across its weight, fill, and optical size axes.", renderIcons},
	{"menu", "Menu", "styles/menu.css", "The list container used for grouped links and settings rows.", renderMenu},
	{"meter", "Meter", "styles/meter.css", "The native meter element, themed across its optimum and sub-optimum ranges.", renderMeter},
	{"pagedlist", "Paged list", "components/pagedlist", "Pagination and debounced search over a placeholder dataset. Typing and paging swap only the list region via htmx.", renderPagedList},
	{"sidebar", "Sidebar", "components/sidebar", "The desktop navigation rail, with grouped items and an active state.", renderSidebar},
	{"bottomtabs", "Bottom tabs", "components/bottomtabs", "The mobile tab bar. Normally only visible below the lg breakpoint.", renderBottomTabs},
	{"bottomsheet", "Bottom sheet", "components/bottomsheet", "A modal sheet anchored to the bottom of the viewport.", renderBottomSheet},
}

func lookup(slug string) (demo, bool) {
	for _, d := range demos {
		if d.Slug == slug {
			return d, true
		}
	}
	return demo{}, false
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

	content, err := current.Render(r)
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

type indexOptions struct {
	Items []template.HTML
}

func renderIndex() (template.HTML, error) {
	items := make([]template.HTML, len(demos))
	for i, d := range demos {
		item, err := indexItemTpl.Render(d)
		if err != nil {
			return "", fmt.Errorf("showcase index: render item %q: %w", d.Slug, err)
		}
		items[i] = item
	}

	content, err := indexTpl.Render(indexOptions{Items: items})
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

func renderDemoPage(current demo, content template.HTML) (template.HTML, error) {
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
