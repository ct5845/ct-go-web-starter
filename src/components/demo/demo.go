package demo

import (
	"html/template"
	"net/url"
)

// Request is what a demo is given when it renders. It carries the page's own
// URL so an interactive demo can build links back to itself without knowing
// where the showcase feature mounted it.
type Request struct {
	BaseHref string
	Query    url.Values
}

// Page is one showcase page. It is declared next to the component it
// demonstrates, so a component and its demo change in the same diff. The
// showcase feature collects these into its index and routes; nothing here has
// an HTTP surface of its own.
type Page struct {
	// Slug is the URL segment, and matches the component's directory name so
	// the page tells you which directory to open.
	Slug string

	Title string

	// Group collects related pages under one primary tab (e.g. "Inputs"),
	// with a secondary tab row listing the group's own pages. Pages with no
	// Group get their own primary tab.
	Group string

	// Source is the path under src/, shown on the page.
	Source string

	Description string

	// Most demos are static and ignore the request.
	Render func(Request) (template.HTML, error)
}
