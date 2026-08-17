package bottomsheet

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/demo"
	_ "embed"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "bottomsheet",
	Title:       "Bottom sheet",
	Source:      "components/bottomsheet",
	Description: "A modal sheet anchored to the bottom of the viewport.",
	Render:      renderShowcase,
}

const showcaseSheetID = "showcase-bottom-sheet"

var (
	//go:embed showcasecontent.html
	showcaseContentHTML string
	showcaseContentTpl  = component.New("showcasecontent.html", showcaseContentHTML)
)

func renderShowcase(demo.Request) (template.HTML, error) {
	content, err := showcaseContentTpl.Render(nil)
	if err != nil {
		return "", err
	}

	sheet, err := Render(Options{
		Id:      showcaseSheetID,
		Content: content,
		Label:   "Showcase bottom sheet",
	})
	if err != nil {
		return "", err
	}

	trigger, err := demo.Frame(demo.FrameOptions{
		Note:    "Uses the native dialog element with the command and commandfor attributes — no JavaScript.",
		Content: template.HTML(`<button class="btn btn-primary" command="show-modal" commandfor="` + showcaseSheetID + `">Open bottom sheet</button>`),
	})
	if err != nil {
		return "", err
	}

	return trigger + sheet, nil
}
