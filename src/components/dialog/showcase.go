package dialog

import (
	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/demo"
	_ "embed"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "dialog",
	Title:       "Dialog",
	Source:      "components/dialog",
	Description: "A modal dialog centered in the viewport.",
	Render:      renderShowcase,
}

const showcaseDialogID = "showcase-dialog"

var (
	//go:embed showcasecontent.html
	showcaseContentHTML string
	showcaseContentTpl  = component.New("showcasecontent.html", showcaseContentHTML)
)

func renderShowcase(demo.Request) (template.HTML, error) {
	content, err := showcaseContentTpl.Render(struct{ Id string }{showcaseDialogID})
	if err != nil {
		return "", err
	}

	sheet, err := Render(Options{
		Id:      showcaseDialogID,
		Content: content,
		Label:   "Showcase dialog",
	})
	if err != nil {
		return "", err
	}

	trigger, err := demo.Frame(demo.FrameOptions{
		Note:    "Uses the native dialog element with the command and commandfor attributes — no JavaScript.",
		Content: template.HTML(`<button class="btn btn-primary" command="show-modal" commandfor="` + showcaseDialogID + `">Open dialog</button>`),
	})
	if err != nil {
		return "", err
	}

	return trigger + sheet, nil
}
