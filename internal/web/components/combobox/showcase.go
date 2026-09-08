package combobox

import (
	"ct-go-web-starter/internal/web/components/demo"
	"html/template"
	"strings"
)

var Showcase = demo.Page{
	Slug:        "combobox",
	Title:       "Combobox",
	Group:       "Inputs",
	Source:      "components/combobox",
	Description: "A searchable, server-driven picker for option sets too large to list up front, opened as a popover anchored to a dropdown-styled trigger.",
	Render:      renderShowcase,
}

var showcaseCountries = []string{
	"Argentina", "Australia", "Belgium", "Brazil", "Canada", "Chile", "China",
	"Denmark", "Egypt", "Finland", "France", "Germany", "Greece", "Hungary",
	"Iceland", "India", "Indonesia", "Ireland", "Israel", "Italy", "Japan",
	"Kenya", "Malaysia", "Mexico", "Morocco", "Netherlands", "New Zealand",
	"Nigeria", "Norway", "Peru", "Philippines", "Poland", "Portugal",
	"Singapore", "South Africa", "South Korea", "Spain", "Sweden",
	"Switzerland", "Thailand", "Turkey", "Uganda", "Ukraine",
	"United Arab Emirates", "United Kingdom", "United States", "Uruguay",
	"Vietnam", "Zambia",
}

func matchingCountries(query string) []string {
	if query == "" {
		return showcaseCountries
	}

	needle := strings.ToLower(query)
	matches := make([]string, 0, len(showcaseCountries))
	for _, country := range showcaseCountries {
		if strings.Contains(strings.ToLower(country), needle) {
			matches = append(matches, country)
		}
	}
	return matches
}

const showcasePageSize = 20

// countryPage returns every match through the given page, so the showcase
// demonstrates combobox's own scroll-to-load-more behaviour rather than
// shipping every match up front. The result accumulates from page 1 rather
// than slicing out just this page, because the "load next page" request
// swaps the whole results wrapper via outerHTML — it needs every row loaded
// so far, not only the newest batch.
func countryPage(query string, page int) (countries []string, hasMore bool) {
	matches := matchingCountries(query)
	end := min(page*showcasePageSize, len(matches))
	return matches[:end], end < len(matches)
}

func countryOption(country string, idPrefix string) Option {
	return Option{
		Id:    idPrefix + strings.ToLower(strings.ReplaceAll(country, " ", "-")),
		Value: country,
		Label: country,
	}
}

func countryOptions(countries []string, idPrefix string) []Option {
	options := make([]Option, len(countries))
	for i, country := range countries {
		options[i] = countryOption(country, idPrefix)
	}
	return options
}

func renderShowcase(request demo.Request) (template.HTML, error) {
	query := request.Query.Get("q")
	page := Page(request.Query)

	countries, hasMore := countryPage(query, page)

	single, err := Render(Options{
		Id:          "showcase-combobox-country",
		Name:        "country",
		Label:       "Country",
		Placeholder: "Choose a country",
		Hint:        "Type to search — 48 countries in this list, 20 loaded at a time.",
		Query:       query,
		Page:        page,
		HasMore:     hasMore,
		BaseHref:    request.BaseHref,
		Selected:    []Option{countryOption("United Kingdom", "")},
		Results:     countryOptions(countries, ""),
	})
	if err != nil {
		return "", err
	}

	multi, err := Render(Options{
		Id:          "showcase-combobox-visited",
		Name:        "visited",
		Label:       "Visited countries",
		Placeholder: "Choose one or more",
		Multi:       true,
		Query:       query,
		Page:        page,
		HasMore:     hasMore,
		BaseHref:    request.BaseHref,
		Selected:    []Option{countryOption("Japan", "multi-"), countryOption("Peru", "multi-")},
		Results:     countryOptions(countries, "multi-"),
	})
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Note:    "Both share the same search box wiring as the paged list: typing debounces an htmx request that swaps just the option list inside the open popover, matching on this showcase's fixed 48-country list server-side. Scrolling to the bottom loads the next 20 matches the same way. A selected item not yet on the loaded page still submits via a hidden input, so the trigger and the form stay correct even before its page loads. Country is single-select (radio); Visited is multi-select (checkbox) and shows every pick in the trigger.",
		Class:   "max-w-md flex flex-col gap-6",
		Content: single + multi,
	})
}
