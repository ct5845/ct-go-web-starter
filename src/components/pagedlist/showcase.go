package pagedlist

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/demo"
	_ "embed"
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

var Showcase = demo.Page{
	Slug:        "pagedlist",
	Title:       "Paged list",
	Source:      "components/pagedlist",
	Description: "Pagination and debounced search over a placeholder dataset. Typing and paging swap only the list region via htmx.",
	Render:      renderShowcase,
}

const showcasePageSize = 10

var (
	//go:embed showcaseentry.html
	showcaseEntryHTML string
	showcaseEntryTpl  = component.New("showcaseentry.html", showcaseEntryHTML)
)

type showcaseEntry struct {
	Label  string
	Detail string
}

var showcaseEntries = buildShowcaseEntries()

func buildShowcaseEntries() []showcaseEntry {
	words := []string{
		"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf",
		"Hotel", "India", "Juliett", "Kilo", "Lima", "Mike", "November",
		"Oscar", "Papa", "Quebec", "Romeo", "Sierra", "Tango", "Uniform",
		"Victor", "Whiskey", "Xray", "Yankee", "Zulu",
	}

	built := make([]showcaseEntry, 0, len(words)*5)
	for round := 1; round <= 5; round++ {
		for _, word := range words {
			built = append(built, showcaseEntry{
				Label:  fmt.Sprintf("%s %d", word, round),
				Detail: fmt.Sprintf("#%03d", len(built)+1),
			})
		}
	}
	return built
}

func matchingShowcaseEntries(query string) []showcaseEntry {
	if query == "" {
		return showcaseEntries
	}

	needle := strings.ToLower(query)
	matches := make([]showcaseEntry, 0, len(showcaseEntries))
	for _, e := range showcaseEntries {
		if strings.Contains(strings.ToLower(e.Label), needle) {
			matches = append(matches, e)
		}
	}
	return matches
}

func renderShowcase(request demo.Request) (template.HTML, error) {
	search := request.Query.Get("q")
	matches := matchingShowcaseEntries(search)
	totalPages := max(1, (len(matches)+showcasePageSize-1)/showcasePageSize)
	pageNumber := showcasePage(request.Query.Get("page"), totalPages)

	start := (pageNumber - 1) * showcasePageSize
	end := min(start+showcasePageSize, len(matches))

	items := make([]template.HTML, 0, end-start)
	for _, match := range matches[start:end] {
		item, err := showcaseEntryTpl.Render(match)
		if err != nil {
			return "", fmt.Errorf("pagedlist showcase: render entry %q: %w", match.Label, err)
		}
		items = append(items, item)
	}

	return Render(Options{
		Items:      items,
		Page:       pageNumber,
		TotalPages: totalPages,
		BaseHref:   request.BaseHref,
		Search: Search{
			Enabled:     true,
			Query:       search,
			Placeholder: "Search entries",
		},
	})
}

// showcasePage clamps the "page" parameter into [1, totalPages] so a stale or
// hand-edited URL cannot slice past the end of the results — easy to hit by
// searching while already on a later page.
func showcasePage(raw string, totalPages int) int {
	pageNumber, err := strconv.Atoi(raw)
	if err != nil {
		return 1
	}
	return min(max(pageNumber, 1), totalPages)
}
