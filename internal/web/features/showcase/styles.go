package showcase

import (
	"ct-go-web-starter/internal/web/components/component"
	"ct-go-web-starter/internal/web/components/demo"
	_ "embed"
	"html/template"
	"time"
)

// These demos document stylesheets rather than components, so unlike the
// component demos there is no package to sit alongside and they live here.

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

	//go:embed menulist.html
	menuListHTML string
	menuListTpl  = component.New("menulist.html", menuListHTML)

	//go:embed meter.html
	meterHTML string
	meterTpl  = component.New("meter.html", meterHTML)

	//go:embed htmxindicator.html
	htmxIndicatorHTML string
	htmxIndicatorTpl  = component.New("htmxindicator.html", htmxIndicatorHTML)
)

var typographyPage = demo.Page{
	Slug:        "typography",
	Title:       "Typography",
	Source:      "static/styles/typography.css",
	Description: "The fluid type scale and base element styles. Resize the window — every size interpolates between its minimum and maximum viewport value.",
	Render: func(demo.Request) (template.HTML, error) {
		return typographyTpl.Render(nil)
	},
}

var colorsPage = demo.Page{
	Slug:        "colors",
	Title:       "Colours",
	Source:      "static/styles/theme.css",
	Description: "Every theme token as a background/foreground pair. Each swatch shows the class names to use.",
	Render: func(demo.Request) (template.HTML, error) {
		return colorsTpl.Render(struct{ Swatches []swatch }{swatches})
	},
}

var spacingPage = demo.Page{
	Slug:        "spacing",
	Title:       "Spacing",
	Source:      "static/styles/spacing.css",
	Description: "Clamped spacing utilities that breathe with the viewport, alongside their fixed equivalents.",
	Render: func(demo.Request) (template.HTML, error) {
		return spacingTpl.Render(nil)
	},
}

var buttonsPage = demo.Page{
	Slug:        "buttons",
	Title:       "Buttons",
	Source:      "static/styles/button.css",
	Description: "Button variants, icon buttons, and disabled states.",
	Render: func(demo.Request) (template.HTML, error) {
		return buttonsTpl.Render(nil)
	},
}

var menuListPage = demo.Page{
	Slug:        "menu-list",
	Title:       "Menu list",
	Source:      "static/styles/menu.css",
	Description: "The .menu/.menu-item list container used for grouped links and settings rows outside of a popover.",
	Render: func(demo.Request) (template.HTML, error) {
		return menuListTpl.Render(nil)
	},
}

var meterPage = demo.Page{
	Slug:        "meter",
	Title:       "Meter",
	Source:      "static/styles/meter.css",
	Description: "The native meter element, themed across its optimum and sub-optimum ranges.",
	Render: func(demo.Request) (template.HTML, error) {
		return meterTpl.Render(nil)
	},
}

// htmxIndicatorDelay slows this demo's own response just enough to see the
// placeholder pulse rather than snap straight to complete on a fast local
// response.
const htmxIndicatorDelay = 1500 * time.Millisecond

type htmxIndicatorResult struct {
	Loaded  bool
	Content template.HTML
}

var htmxIndicatorPage = demo.Page{
	Slug:        "htmx-indicator",
	Title:       "Htmx indicator",
	Source:      "static/styles/htmxindicator.css",
	Description: "The .htmx-indicator utility: hidden until htmx toggles its \"htmx-request\" class on the element an hx-indicator points at, so a swapped-out region can show a placeholder with no JS of its own.",
	Render: func(request demo.Request) (template.HTML, error) {
		// The page itself renders instantly with an empty, unloaded result;
		// its own hx-trigger="load" immediately re-requests with load=true,
		// so every page view (including a plain refresh) shows the
		// placeholder for htmxIndicatorDelay before the real content lands.
		// The loaded response drops hx-trigger="load" from the swapped-in
		// markup, since re-adding it would re-fire on every swap.
		if request.Query.Get("load") != "true" {
			return htmxIndicatorTpl.Render(htmxIndicatorResult{})
		}

		time.Sleep(htmxIndicatorDelay)
		return htmxIndicatorTpl.Render(htmxIndicatorResult{
			Loaded:  true,
			Content: template.HTML("Loaded at " + time.Now().Format("15:04:05.000")),
		})
	},
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
