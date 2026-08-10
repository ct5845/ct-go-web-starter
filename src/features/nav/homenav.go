package nav

import (
	"ct-go-web-starter/src/components/bottomtabs"
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/sidebar"
	_ "embed"
	"html/template"
)

var (
	//go:embed homenav.html
	homenavHTML string
	homenavTpl  = component.New("bottomsheetnav.html", homenavHTML)
)

const moreSheetID = "home-more-sheet"

type primaryItem struct {
	Label string
	Href  string
	Icon  string
	Tab   string
}

// primaryItems is the single source of truth for the main navigation. Both the
// mobile bottom tabs and the desktop sidebar are projected from it, so the two
// cannot drift apart.
var primaryItems = []primaryItem{
	{Label: "Home", Href: "/", Icon: "home", Tab: "home"},
	{Label: "Option A", Href: "/option-a", Icon: "colors", Tab: "option-a"},
}

func createTabs(activeTab string) (template.HTML, error) {
	tabs := make([]bottomtabs.Tab, 0, len(primaryItems)+1)
	for _, it := range primaryItems {
		tabs = append(tabs, bottomtabs.Tab{Label: it.Label, Href: it.Href, Active: activeTab == it.Tab, Icon: it.Icon})
	}
	tabs = append(tabs, bottomtabs.Tab{Label: "More", Icon: "menu", Attrs: `command="show-modal" commandfor="` + moreSheetID + `"`})
	return bottomtabs.Render(bottomtabs.Options{Tabs: tabs})
}

func createSidebar(activeTab string) (template.HTML, error) {
	header, err := RenderSidebarHeader()
	if err != nil {
		return "", err
	}

	items := make([]sidebar.Item, len(primaryItems))
	for i, it := range primaryItems {
		items[i] = sidebar.Item{Label: it.Label, Href: it.Href, Active: activeTab == it.Tab, Icon: it.Icon}
	}

	return sidebar.Render(sidebar.Options{
		Header: header,
		Groups: []sidebar.Group{{Items: items}},
	})
}

// Result holds the navigation rendered for both layouts: Footer is the mobile
// bottom tabs plus the "More" sheet, SideBar is the desktop sidebar. The layout
// places each in the region a media query reveals.
type Result struct {
	Footer  template.HTML
	SideBar template.HTML
}

func Render(activeTab string) (Result, error) {
	bottomSheet, err := createBottomSheet(activeTab)
	if err != nil {
		return Result{}, err
	}

	bottomTabs, err := createTabs(activeTab)
	if err != nil {
		return Result{}, err
	}

	footer, err := homenavTpl.Render(struct {
		BottomTabs  template.HTML
		BottomSheet template.HTML
	}{
		BottomTabs:  bottomTabs,
		BottomSheet: bottomSheet,
	})
	if err != nil {
		return Result{}, err
	}

	sideBar, err := createSidebar(activeTab)
	if err != nil {
		return Result{}, err
	}

	return Result{Footer: footer, SideBar: sideBar}, nil
}
