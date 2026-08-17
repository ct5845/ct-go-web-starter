package icon

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/demo"
	_ "embed"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "icon",
	Title:       "Icons",
	Source:      "components/icon",
	Description: "The subsetted icon font across its weight, fill, and optical size axes.",
	Render:      renderShowcase,
}

var (
	//go:embed showcase.html
	showcaseHTML string
	showcaseTpl  = component.New("showcase.html", showcaseHTML)
)

// variant pairs a set of icon utility classes with a label, so the axes of the
// variable font can be compared side by side.
type variant struct {
	Class string
	Label string
}

var variants = []variant{
	{"icon-wght-200", "wght 200"},
	{"icon-wght-400", "wght 400"},
	{"icon-wght-700", "wght 700"},
	{"icon-fill-1", "fill 1"},
	{"icon-opsz-20", "opsz 20"},
	{"icon-opsz-40", "opsz 40"},
	{"icon-opsz-48", "opsz 48"},
}

func renderShowcase(demo.Request) (template.HTML, error) {
	return showcaseTpl.Render(struct {
		Names    []string
		Variants []variant
	}{Names, variants})
}
