package combobox

import (
	"ct-go-web-starter/internal/web/components/component"
	"ct-go-web-starter/internal/web/components/field"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

var (
	//go:embed combobox.html
	comboboxHTML string
	//go:embed combobox.js
	comboboxJS  string
	comboboxTpl = component.WithAlpine("combobox.html", comboboxHTML, comboboxJS)

	//go:embed option.html
	optionHTML string
	optionTpl  = component.New("option.html", optionHTML)

	//go:embed hidden-option.html
	hiddenOptionHTML string
	hiddenOptionTpl  = component.New("hidden-option.html", hiddenOptionHTML)
)

// Option is one search result. Id must be unique across the page, since it
// backs the hidden radio/checkbox input and its label's for attribute.
// Whether it's checked is decided by membership in Options.Selected, not a
// field here, so a value can never disagree with itself between Results and
// Selected.
type Option struct {
	Id    string
	Value string
	Label string
}

// Options configures a searchable, server-driven picker for option sets too
// large to list up front (e.g. one of a few hundred countries), opened as a
// popover anchored to a dropdown-styled trigger. Multi switches between
// single (radio) and multi (checkbox) selection; either way the popover's
// search box requests BaseHref with the query as "q" and expects the same
// Options shape back, so the caller re-renders Options.Results directly as
// an htmx partial for that request — see the showcase for the full
// request/response cycle.
//
// Results should be capped to one page (e.g. 20) rather than the whole
// dataset — HasMore drives a scroll-triggered request for the next Page,
// appended onto the open list, so scrolling behaves like search: a request
// per page rather than shipping everything up front.
type Options struct {
	Id          string
	Name        string
	Label       string
	Placeholder string
	Multi       bool

	// Selected is the full, authoritative set of chosen items, independent
	// of paging or search — the caller must supply these on every render
	// (e.g. by looking them up by id), since a selection made on an earlier
	// page won't otherwise appear once search or scroll has moved past it.
	// Any entry not already present in Results submits via a hidden input
	// instead (see hidden-option.html), so the form always carries the full
	// selection even before the matching page loads.
	Selected []Option

	Results     []Option
	Query       string
	Page        int
	HasMore     bool
	BaseHref    string
	ResultsNote string

	Hint     string
	Error    string
	Required bool
}

// ResultsId is the DOM id of just the swappable option-list region inside
// the popover, for a handler to target with hx-select/hx-target when
// rendering Options.Results alone in response to a search request — see
// IsSearchRequest.
func ResultsId(id string) string {
	return id + "-results"
}

// IsSearchRequest reports whether r is an htmx search request from a
// combobox's popover, so the caller can respond with a Results-only partial
// instead of a full page.
func IsSearchRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// Page parses the "page" query parameter a scroll or search request carries,
// defaulting to 1 for the initial render or a malformed value.
func Page(query url.Values) int {
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func pageHref(baseHref, query string, page int) string {
	u, err := url.Parse(baseHref)
	if err != nil {
		return fmt.Sprintf("%s?q=%s&page=%d", baseHref, url.QueryEscape(query), page)
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return u.String()
}

type templateOptions struct {
	Id            string
	Name          string
	Label         string
	Placeholder   string
	Multi         bool
	InputType     string
	SelectedText  string
	ResultsId     string
	HiddenOptions []template.HTML
	Results       []template.HTML
	Query         string
	BaseHref      string
	NextPageHref  string
	ResultsNote   string
	Required      bool
	DescribedBy   template.HTMLAttr
	Invalid       template.HTMLAttr
}

func Render(options Options) (template.HTML, error) {
	described := field.Describe(options.Id, options.Hint, options.Error)

	inputType := "radio"
	if options.Multi {
		inputType = "checkbox"
	}

	selectedValues := make(map[string]bool, len(options.Selected))
	onPage := make(map[string]bool, len(options.Results))
	for _, o := range options.Selected {
		selectedValues[o.Value] = true
	}
	for _, o := range options.Results {
		onPage[o.Value] = true
	}

	// A selected item not on the loaded page still needs to submit, so it
	// gets a hidden input outside the visible list rather than a styled
	// row — see hidden-option.html. Once the real row loads (matched by
	// value in combobox.js), the hidden one is removed.
	var hiddenOptions []template.HTML
	for _, o := range options.Selected {
		if onPage[o.Value] {
			continue
		}
		rendered, err := o.renderHidden(options.Id, options.Name, inputType)
		if err != nil {
			return "", err
		}
		hiddenOptions = append(hiddenOptions, rendered)
	}

	results := make([]template.HTML, len(options.Results))
	for i, o := range options.Results {
		rendered, err := o.render(options.Id, options.Name, inputType, selectedValues[o.Value])
		if err != nil {
			return "", err
		}
		results[i] = rendered
	}

	selectedText := options.Placeholder
	if len(options.Selected) > 0 {
		selectedText = options.Selected[0].Label
		for _, o := range options.Selected[1:] {
			selectedText += ", " + o.Label
		}
	}

	var nextPageHref string
	if options.HasMore {
		page := options.Page
		if page < 1 {
			page = 1
		}
		nextPageHref = pageHref(options.BaseHref, options.Query, page+1)
	}

	input, err := comboboxTpl.Render(templateOptions{
		Id:            options.Id,
		Name:          options.Name,
		Label:         options.Label,
		Placeholder:   options.Placeholder,
		Multi:         options.Multi,
		InputType:     inputType,
		SelectedText:  selectedText,
		ResultsId:     ResultsId(options.Id),
		HiddenOptions: hiddenOptions,
		Results:       results,
		Query:         options.Query,
		BaseHref:      options.BaseHref,
		NextPageHref:  nextPageHref,
		ResultsNote:   options.ResultsNote,
		Required:      options.Required,
		DescribedBy:   described.DescribedBy,
		Invalid:       described.Invalid,
	})
	if err != nil {
		return "", err
	}

	return field.Render(field.Options{
		Id:       options.Id,
		Label:    options.Label,
		Hint:     options.Hint,
		Error:    options.Error,
		Required: options.Required,
		Input:    input,
		NoLabel:  true,
	})
}

type optionProps struct {
	Id      string
	Name    string
	Type    string
	Value   string
	Label   string
	Checked bool
}

func (o Option) render(comboboxId, name, inputType string, checked bool) (template.HTML, error) {
	return optionTpl.Render(optionProps{
		Id:      comboboxId + "-option-" + o.Id,
		Name:    name,
		Type:    inputType,
		Value:   o.Value,
		Label:   o.Label,
		Checked: checked,
	})
}

// renderHidden submits a selection that isn't on the loaded page yet — see
// hidden-option.html. combobox.js matches it to its eventual visible row by
// value (its id is distinct, since both can briefly coexist in the DOM
// until combobox.js removes this one once that row loads).
func (o Option) renderHidden(comboboxId, name, inputType string) (template.HTML, error) {
	return hiddenOptionTpl.Render(optionProps{
		Id:    comboboxId + "-hidden-option-" + o.Id,
		Name:  name,
		Type:  inputType,
		Value: o.Value,
		Label: o.Label,
	})
}
