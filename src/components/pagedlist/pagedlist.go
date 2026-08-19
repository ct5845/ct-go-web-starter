package pagedlist

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var p = message.NewPrinter(language.English)

//go:embed pagedlist.html
var pagedlistHTML string

//go:embed pagedlist.js
var pagedlistJS string

var pagedlistTpl = component.WithIIFE("pagedlist.html", pagedlistHTML, pagedlistJS)

// id is the DOM id of the whole component; resultsID is the id of just the
// swappable list+pagination region, excluding the search input. Only one
// paged list is expected per page, so fixed ids are enough.
//
// The input must live outside resultsID: htmx replaces whatever hx-select
// matches on every keystroke, and replacing the input itself would drop
// focus and cursor position mid-typing.
const (
	id        = "pagedlist"
	resultsID = "pagedlist-results"
)

type Options struct {
	Items []template.HTML

	// Page is 1-indexed.
	Page int

	TotalPages int

	// BaseHref is the path the command bar and search box request against;
	// "page" and "q" are added or replaced as query parameters. It may
	// already carry other query parameters (e.g. "/admin/domains?sort=count").
	BaseHref string

	// Search enables the search box above the list. Zero value leaves it out.
	Search Search

	Theme string
}

type Search struct {
	// Enabled turns the search box on. When off, Query/Placeholder are unused.
	Enabled bool

	Query string

	Placeholder string
}

type templateOptions struct {
	ID            string
	ResultsID     string
	Items         []template.HTML
	Page          string
	TotalPages    string
	FirstHref     string
	PrevHref      string
	NextHref      string
	LastHref      string
	HasPrev       bool
	HasNext       bool
	SearchEnabled bool
	SearchQuery   string
	Placeholder   string
	BaseHref      string
	Theme         string
}

func Render(options Options) (template.HTML, error) {
	base := options.BaseHref
	if options.Search.Query != "" {
		base = withQuery(base, options.Search.Query)
	}
	theme := "bg-secondary-container text-on-secondary-container"
	if options.Theme != "" {
		theme = options.Theme
	}

	return pagedlistTpl.Render(templateOptions{
		ID:            id,
		ResultsID:     resultsID,
		Items:         options.Items,
		Page:          p.Sprintf("%d", options.Page),
		TotalPages:    p.Sprintf("%d", options.TotalPages),
		FirstHref:     hrefForPage(base, 1),
		PrevHref:      hrefForPage(base, options.Page-1),
		NextHref:      hrefForPage(base, options.Page+1),
		LastHref:      hrefForPage(base, options.TotalPages),
		HasPrev:       options.Page > 1,
		HasNext:       options.Page < options.TotalPages,
		SearchEnabled: options.Search.Enabled,
		SearchQuery:   options.Search.Query,
		Placeholder:   options.Search.Placeholder,
		BaseHref:      options.BaseHref,
		Theme:         theme,
	})
}

func withQuery(baseHref, query string) string {
	u, err := url.Parse(baseHref)
	if err != nil {
		return fmt.Sprintf("%s?q=%s", baseHref, url.QueryEscape(query))
	}
	q := u.Query()
	q.Set("q", query)
	u.RawQuery = q.Encode()
	return u.String()
}

func hrefForPage(baseHref string, page int) string {
	u, err := url.Parse(baseHref)
	if err != nil {
		return fmt.Sprintf("%s?page=%d", baseHref, page)
	}
	q := u.Query()
	q.Set("page", fmt.Sprint(page))
	u.RawQuery = q.Encode()
	return u.String()
}

// IsPartialRequest reports whether r is an htmx request targeting this
// component, so the caller can respond with just the Render output instead of
// a full page.
func IsPartialRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
