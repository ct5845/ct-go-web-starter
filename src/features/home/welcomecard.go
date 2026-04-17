package home

import (
	_ "embed"
	"html/template"

	"ct-go-web-starter/src/components/component"
)

//go:embed welcomecard.html
var welcomeCardHTML string
var welcomeCardTmpl = component.New("welcomecard.html", welcomeCardHTML)

func renderWelcomeCard(title, description string) (template.HTML, error) {
	return welcomeCardTmpl.Render(
		"Title", title,
		"Description", description,
	)
}
