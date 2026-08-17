package showcase

import (
	"ct-go-web-starter/src/components/bottomsheet"
	"ct-go-web-starter/src/components/bottomtabs"
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/icon"
	"ct-go-web-starter/src/components/sidebar"
	_ "embed"
	"html/template"
	"net/http"
)

var (
	//go:embed typography.html
	typographyHTML string
	typographyTpl  = component.New("typography.html", typographyHTML)

	//go:embed colors.html
	colorsHTML string
	colorsTpl  = component.New("colors.html", colorsHTML)

	//go:embed spacing.html
	spacingHTML string
	spacingTpl  = component.New("spacing.html", spacingHTML)

	//go:embed buttons.html
	buttonsHTML string
	buttonsTpl  = component.New("buttons.html", buttonsHTML)

	//go:embed icons.html
	iconsHTML string
	iconsTpl  = component.New("icons.html", iconsHTML)

	//go:embed menu.html
	menuHTML string
	menuTpl  = component.New("menu.html", menuHTML)

	//go:embed meter.html
	meterHTML string
	meterTpl  = component.New("meter.html", meterHTML)

	//go:embed frame.html
	frameHTML string
	frameTpl  = component.New("frame.html", frameHTML)
)

func renderTypography(*http.Request) (template.HTML, error) {
	return typographyTpl.Render(nil)
}

func renderSpacing(*http.Request) (template.HTML, error) {
	return spacingTpl.Render(nil)
}

func renderButtons(*http.Request) (template.HTML, error) {
	return buttonsTpl.Render(nil)
}

func renderMenu(*http.Request) (template.HTML, error) {
	return menuTpl.Render(nil)
}

func renderMeter(*http.Request) (template.HTML, error) {
	return meterTpl.Render(nil)
}

// swatch is one theme colour rendered as a background/foreground pair, so the
// pairing is visibly legible rather than only documented.
type swatch struct {
	Background string
	Foreground string
	Label      string
}

var swatches = []swatch{
	{"bg-primary", "text-on-primary", "primary"},
	{"bg-primary-container", "text-on-primary-container", "primary-container"},
	{"bg-secondary", "text-on-secondary", "secondary"},
	{"bg-secondary-container", "text-on-secondary-container", "secondary-container"},
	{"bg-tertiary", "text-on-tertiary", "tertiary"},
	{"bg-tertiary-container", "text-on-tertiary-container", "tertiary-container"},
	{"bg-error", "text-on-error", "error"},
	{"bg-error-container", "text-on-error-container", "error-container"},
	{"bg-surface", "text-on-surface", "surface"},
	{"bg-surface-variant", "text-on-surface-variant", "surface-variant"},
	{"bg-surface-dim", "text-on-surface", "surface-dim"},
	{"bg-surface-bright", "text-on-surface", "surface-bright"},
	{"bg-inverse-surface", "text-on-inverse-surface", "inverse-surface"},
	{"bg-container-lowest", "text-on-surface", "container-lowest"},
	{"bg-container-low", "text-on-surface", "container-low"},
	{"bg-container", "text-on-surface", "container"},
	{"bg-container-high", "text-on-surface", "container-high"},
	{"bg-container-highest", "text-on-surface", "container-highest"},
	{"bg-outline", "text-on-inverse-surface", "outline"},
	{"bg-outline-variant", "text-on-surface", "outline-variant"},
}

func renderColors(*http.Request) (template.HTML, error) {
	return colorsTpl.Render(struct{ Swatches []swatch }{swatches})
}

// iconVariant pairs a set of icon utility classes with the label describing it,
// so the axes of the variable icon font can be compared side by side.
type iconVariant struct {
	Class string
	Label string
}

var iconVariants = []iconVariant{
	{"icon-wght-200", "wght 200"},
	{"icon-wght-400", "wght 400"},
	{"icon-wght-700", "wght 700"},
	{"icon-fill-1", "fill 1"},
	{"icon-opsz-20", "opsz 20"},
	{"icon-opsz-40", "opsz 40"},
	{"icon-opsz-48", "opsz 48"},
}

func renderIcons(*http.Request) (template.HTML, error) {
	return iconsTpl.Render(struct {
		Names    []string
		Variants []iconVariant
	}{icon.Names, iconVariants})
}

func renderSidebar(*http.Request) (template.HTML, error) {
	rendered, err := sidebar.Render(sidebar.Options{
		Groups: []sidebar.Group{
			{Items: []sidebar.Item{
				{Label: "Home", Href: "#", Icon: "home", Active: true},
				{Label: "Search", Href: "#", Icon: "search"},
			}},
			{Title: "Group", Items: []sidebar.Item{
				{Label: "Showcase", Href: "#", Icon: "widgets"},
				{Label: "Settings", Href: "#", Icon: "colors"},
			}},
		},
	})
	if err != nil {
		return "", err
	}

	return frameTpl.Render(frameProps{
		Note:    "Rendered at its natural width. In a page this sits in the aside of layoutswitch, visible from the lg breakpoint up.",
		Class:   "w-64 border border-outline-variant rounded-xl overflow-hidden",
		Content: rendered,
	})
}

func renderBottomTabs(*http.Request) (template.HTML, error) {
	rendered, err := bottomtabs.Render(bottomtabs.Options{
		Tabs: []bottomtabs.Tab{
			{Label: "Home", Href: "#", Icon: "home", Active: true},
			{Label: "Search", Href: "#", Icon: "search"},
			{Label: "More", Icon: "menu"},
		},
	})
	if err != nil {
		return "", err
	}

	return frameTpl.Render(frameProps{
		Note:    "Normally pinned to the bottom of the viewport below the lg breakpoint. Shown inline here.",
		Class:   "max-w-md border border-outline-variant rounded-xl overflow-hidden",
		Content: rendered,
	})
}

const showcaseSheetID = "showcase-sheet"

var (
	//go:embed bottomsheet.html
	bottomsheetHTML string
	bottomsheetTpl  = component.New("bottomsheet.html", bottomsheetHTML)
)

func renderBottomSheet(*http.Request) (template.HTML, error) {
	content, err := bottomsheetTpl.Render(nil)
	if err != nil {
		return "", err
	}

	sheet, err := bottomsheet.Render(bottomsheet.Options{
		Id:      showcaseSheetID,
		Content: content,
		Label:   "Showcase bottom sheet",
	})
	if err != nil {
		return "", err
	}

	trigger, err := frameTpl.Render(frameProps{
		Note:    "Uses the native dialog element with the command/commandfor attributes — no JavaScript.",
		Class:   "",
		Content: template.HTML(`<button class="btn btn-primary" command="show-modal" commandfor="` + showcaseSheetID + `">Open bottom sheet</button>`),
	})
	if err != nil {
		return "", err
	}

	return trigger + sheet, nil
}

type frameProps struct {
	Note    string
	Class   string
	Content template.HTML
}
