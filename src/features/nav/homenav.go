package nav

import (
	"ct-go-web-starter/src/components/bottomtabs"
	"ct-go-web-starter/src/components/component"
	_ "embed"
	"html/template"
)

var (
	//go:embed homenav.html
	homenavHTML string
	homenavTpl  = component.New("bottomsheetnav.html", homenavHTML)
)

const moreSheetID = "home-more-sheet"

func createTabs(activeTab string) (template.HTML, error) {
	return bottomtabs.Render(bottomtabs.Options{
		Tabs: []bottomtabs.Tab{
			{Label: "Home", Href: "/", Active: activeTab == "home", Icon: "home"},
			{Label: "Option A", Href: "/option-a", Active: activeTab == "option-a", Icon: "colors"},
			{Label: "More", Icon: "menu", Attrs: `command="show-modal" commandfor="` + moreSheetID + `"`},
		},
	})
}

func Render(activeTab string) (template.HTML, error) {
	bottomSheet, err := createBottomSheet(activeTab)
	if err != nil {
		return "", err
	}

	bottomTabs, err := createTabs(activeTab)
	if err != nil {
		return "", err
	}

	return homenavTpl.Render(struct {
		BottomTabs  template.HTML
		BottomSheet template.HTML
	}{
		BottomTabs:  bottomTabs,
		BottomSheet: bottomSheet,
	})
}
