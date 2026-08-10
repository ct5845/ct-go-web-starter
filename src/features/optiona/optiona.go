package optiona

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/layoutswitch"
	"ct-go-web-starter/src/components/page"
	"ct-go-web-starter/src/components/pagedlist"
	"ct-go-web-starter/src/features/nav"
	"ct-go-web-starter/src/infrastructure/reqlog"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

const pageSize = 10

var (
	//go:embed optiona.html
	optionaHTML string
	optionaTpl  = component.New("optiona.html", optionaHTML)

	//go:embed entry.html
	entryHTML string
	entryTpl  = component.New("entry.html", entryHTML)
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /option-a", handleGet)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	defer reqlog.Track(r.Context(), "optiona.handleGet", "")()

	query := r.URL.Query().Get("q")
	matches := matching(query)
	totalPages := max(1, (len(matches)+pageSize-1)/pageSize)
	pageNumber := requestedPage(r.URL.Query().Get("page"), totalPages)

	list, err := renderList(matches, pageNumber, totalPages, query)
	if err != nil {
		slog.Error("Failed to render option A list", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// A search keystroke or a boosted pagination link only needs the list
	// region back; htmx selects it out of the response itself.
	if pagedlist.IsPartialRequest(r) {
		io.WriteString(w, string(list))
		return
	}

	full, err := renderPage(list, len(matches))
	if err != nil {
		slog.Error("Failed to render option A page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(full))
}

// requestedPage clamps the "page" query parameter into [1, totalPages] so a
// hand-edited or stale URL cannot slice past the end of the results.
func requestedPage(raw string, totalPages int) int {
	pageNumber, err := strconv.Atoi(raw)
	if err != nil {
		return 1
	}
	return min(max(pageNumber, 1), totalPages)
}

func renderList(matches []entry, pageNumber, totalPages int, query string) (template.HTML, error) {
	start := (pageNumber - 1) * pageSize
	end := min(start+pageSize, len(matches))

	items := make([]template.HTML, 0, end-start)
	for _, match := range matches[start:end] {
		item, err := entryTpl.Render(match)
		if err != nil {
			return "", fmt.Errorf("option A: render entry %q: %w", match.Label, err)
		}
		items = append(items, item)
	}

	return pagedlist.Render(pagedlist.Options{
		Items:      items,
		Page:       pageNumber,
		TotalPages: totalPages,
		BaseHref:   "/option-a",
		Search: pagedlist.Search{
			Enabled:     true,
			Query:       query,
			Placeholder: "Search entries",
		},
	})
}

type options struct {
	Title       string
	Description string
	List        template.HTML
}

func renderPage(list template.HTML, matchCount int) (template.HTML, error) {
	content, err := optionaTpl.Render(options{
		Title:       "Option A",
		Description: fmt.Sprintf("A demo of the pagedlist component: %d placeholder entries, %d per page, with debounced search.", matchCount, pageSize),
		List:        list,
	})
	if err != nil {
		return "", fmt.Errorf("option A page: render content: %w", err)
	}

	navigation, err := nav.Render("option-a")
	if err != nil {
		return "", fmt.Errorf("option A page: render navigation: %w", err)
	}

	return layoutswitch.RenderPage(page.Options{
		Title:           "Option A",
		MetaDescription: "A paginated, searchable list built with the pagedlist component",
	}, layoutswitch.Options{
		Content:    content,
		BottomTabs: navigation.Footer,
		SideBar:    navigation.SideBar,
	})
}
