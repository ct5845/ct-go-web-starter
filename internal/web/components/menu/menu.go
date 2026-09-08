package menu

import (
	"ct-go-web-starter/internal/web/components/component"
	_ "embed"
	"fmt"
	"html/template"
)

var (
	//go:embed menu.html
	menuHTML string
	menuTpl  = component.New("menu.html", menuHTML)

	//go:embed item.html
	itemHTML string
	itemTpl  = component.New("item.html", itemHTML)
)

// Item is one row in a menu. A leaf item sets Href; an item with Items renders
// as a submenu trigger that opens a nested popover positioned off itself.
type Item struct {
	Label string
	Icon  string
	Href  string
	Items []Item
}

// Options configures the top-level popover. Id must match the popovertarget
// and anchor-name the consumer sets on the trigger button — see the showcase
// for the expected trigger markup.
type Options struct {
	Id    string
	Items []Item
}

// Render draws a popover menu. The caller owns the trigger button: it needs
// popovertarget="{{.Id}}" and an inline anchor-name matching AnchorName(Id),
// so the popover below can position itself with position-anchor.
func Render(options Options) (template.HTML, error) {
	return renderPopover(options.Id, options.Items, false)
}

// AnchorName is the CSS anchor-name a trigger must declare to position a
// popover with the given id against itself.
func AnchorName(id string) string {
	return "--anchor-" + id
}

type popoverProps struct {
	Id    string
	Style template.CSS
	Items []template.HTML
}

func renderPopover(id string, items []Item, nested bool) (template.HTML, error) {
	rendered := make([]template.HTML, len(items))
	for i, it := range items {
		item, err := it.render(id, i)
		if err != nil {
			return "", err
		}
		rendered[i] = item
	}

	// html/template's CSS sanitizer does not recognise anchor-positioning
	// syntax (dashed-idents, multi-keyword position-area values) as safe and
	// silently blanks it out, so the whole declaration is built as one
	// trusted string rather than composed from template actions.
	var area, tryFallback, marginProperty string
	if nested {
		// A submenu opens beside its trigger item, like a native OS submenu,
		// flipping to the left when it would run off the right edge.
		area, tryFallback, marginProperty = "right span-bottom", "flip-inline", "margin-inline-start"
	} else {
		// The top-level menu opens below its trigger button, flipping above
		// it when it would run off the bottom edge.
		area, tryFallback, marginProperty = "bottom span-right", "flip-block", "margin-block-start"
	}
	style := fmt.Sprintf(
		"position-anchor: %s; position-area: %s; position-try-fallbacks: %s; inset: auto; %s: 0.25rem",
		AnchorName(id), area, tryFallback, marginProperty,
	)

	return menuTpl.Render(popoverProps{
		Id:    id,
		Style: template.CSS(style),
		Items: rendered,
	})
}

type itemProps struct {
	Label     string
	Icon      string
	Href      string
	IsSubmenu bool
	SubmenuId string
	Style     template.CSS
	Submenu   template.HTML
}

func (it Item) render(parentId string, index int) (template.HTML, error) {
	if len(it.Items) == 0 {
		return itemTpl.Render(itemProps{
			Label: it.Label,
			Icon:  it.Icon,
			Href:  it.Href,
		})
	}

	submenuId := fmt.Sprintf("%s-%d", parentId, index)
	submenu, err := renderPopover(submenuId, it.Items, true)
	if err != nil {
		return "", err
	}

	return itemTpl.Render(itemProps{
		Label:     it.Label,
		Icon:      it.Icon,
		IsSubmenu: true,
		SubmenuId: submenuId,
		Style:     template.CSS("anchor-name: " + AnchorName(submenuId)),
		Submenu:   submenu,
	})
}
