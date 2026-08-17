package showcase

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/pagedlist"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const pageSize = 10

var (
	//go:embed pagedlistentry.html
	entryHTML string
	entryTpl  = component.New("pagedlistentry.html", entryHTML)
)

type entry struct {
	Label  string
	Detail string
}

// entries is placeholder data so the paged list has enough rows to page and
// search through.
var entries = buildEntries()

func buildEntries() []entry {
	words := []string{
		"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf",
		"Hotel", "India", "Juliett", "Kilo", "Lima", "Mike", "November",
		"Oscar", "Papa", "Quebec", "Romeo", "Sierra", "Tango", "Uniform",
		"Victor", "Whiskey", "Xray", "Yankee", "Zulu",
	}

	built := make([]entry, 0, len(words)*5)
	for round := 1; round <= 5; round++ {
		for _, word := range words {
			built = append(built, entry{
				Label:  fmt.Sprintf("%s %d", word, round),
				Detail: fmt.Sprintf("#%03d", len(built)+1),
			})
		}
	}
	return built
}

func matching(query string) []entry {
	if query == "" {
		return entries
	}

	needle := strings.ToLower(query)
	matches := make([]entry, 0, len(entries))
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Label), needle) {
			matches = append(matches, e)
		}
	}
	return matches
}

func renderPagedList(r *http.Request) (template.HTML, error) {
	query := r.URL.Query().Get("q")
	matches := matching(query)
	totalPages := max(1, (len(matches)+pageSize-1)/pageSize)
	pageNumber := requestedPage(r.URL.Query().Get("page"), totalPages)

	start := (pageNumber - 1) * pageSize
	end := min(start+pageSize, len(matches))

	items := make([]template.HTML, 0, end-start)
	for _, match := range matches[start:end] {
		item, err := entryTpl.Render(match)
		if err != nil {
			return "", fmt.Errorf("showcase pagedlist: render entry %q: %w", match.Label, err)
		}
		items = append(items, item)
	}

	return pagedlist.Render(pagedlist.Options{
		Items:      items,
		Page:       pageNumber,
		TotalPages: totalPages,
		BaseHref:   "/showcase/pagedlist",
		Search: pagedlist.Search{
			Enabled:     true,
			Query:       query,
			Placeholder: "Search entries",
		},
	})
}

// requestedPage clamps the "page" query parameter into [1, totalPages] so a
// hand-edited or stale URL cannot slice past the end of the results — easy to
// hit by searching while already on a later page.
func requestedPage(raw string, totalPages int) int {
	pageNumber, err := strconv.Atoi(raw)
	if err != nil {
		return 1
	}
	return min(max(pageNumber, 1), totalPages)
}
