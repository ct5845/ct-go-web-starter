package sidebar

import (
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

var (
	//go:embed sidebar.html
	sidebarHTML string
	sidebarTpl  = component.New("sidebar.html", sidebarHTML)

	//go:embed item.html
	itemHTML string
	itemTpl  = component.New("item.html", itemHTML)

	//go:embed group.html
	groupHTML string
	groupTpl  = component.New("group.html", groupHTML)
)

type Group struct {
	Title string
	Items []Item
}

type Item struct {
	Icon   string
	Label  string
	Href   string
	Active bool
}

type Options struct {
	Header template.HTML
	Groups []Group
	Footer template.HTML
}

type templateOptions struct {
	Header template.HTML
	Groups []template.HTML
	Footer template.HTML
}

func Render(options Options) (template.HTML, error) {
	groups := make([]template.HTML, len(options.Groups))
	for i, group := range options.Groups {
		rendered, err := group.render()
		if err != nil {
			return "", err
		}
		groups[i] = rendered
	}
	return sidebarTpl.Render(templateOptions{Header: options.Header, Groups: groups, Footer: options.Footer})
}

type itemProps struct {
	Href        string
	AriaCurrent template.HTMLAttr
	Class       string
	Icon        string
	IconClass   string
	Label       string
}

func (it Item) render() (template.HTML, error) {
	iconClass := "icon-wght-400 icon-opsz-20"
	stateClass := "text-on-surface-variant hover:bg-container"
	var ariaCurrent template.HTMLAttr
	if it.Active {
		iconClass = "icon-fill-1 icon-wght-700 icon-opsz-20"
		stateClass = "text-on-primary-container bg-primary-container font-bold"
		ariaCurrent = `aria-current="page"`
	}

	return itemTpl.Render(itemProps{
		Href:        it.Href,
		AriaCurrent: ariaCurrent,
		Class:       "flex items-center gap-2 rounded-lg p-2 " + stateClass,
		Icon:        it.Icon,
		IconClass:   iconClass,
		Label:       it.Label,
	})
}

type groupProps struct {
	Title string
	Items []template.HTML
}

func (g Group) render() (template.HTML, error) {
	items := make([]template.HTML, len(g.Items))
	for i, item := range g.Items {
		rendered, err := item.render()
		if err != nil {
			return "", err
		}
		items[i] = rendered
	}
	return groupTpl.Render(groupProps{Title: g.Title, Items: items})
}
