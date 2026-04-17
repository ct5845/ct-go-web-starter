package home

import (
	_ "embed"
	"fmt"
	"html/template"

	"ct-go-web-starter/src/components/component"
	"ct-go-web-starter/src/components/footer"
	"ct-go-web-starter/src/components/header"
	"ct-go-web-starter/src/components/page"
)

//go:embed home.html
var homeHTML string
var homeTmpl = component.New("home.html", homeHTML)

func renderPage() (template.HTML, error) {
	headerHTML, err := header.Render("Title", "CT Go Web Starter")
	if err != nil {
		return "", fmt.Errorf("home page: render header: %w", err)
	}

	footerHTML, err := footer.Render()
	if err != nil {
		return "", fmt.Errorf("home page: render footer: %w", err)
	}

	welcomeCardHTML, err := renderWelcomeCard("CT Go Web Starter", "A modern Go web application starter with HTMX, Alpine.js, and TailwindCSS")
	if err != nil {
		return "", fmt.Errorf("home page: render welcome card: %w", err)
	}

	contentHTML, err := homeTmpl.Render("WelcomeCardHTML", welcomeCardHTML)
	if err != nil {
		return "", fmt.Errorf("home page: render content: %w", err)
	}

	pageHTML, err := page.Render(
		"Title", "CT Go Web Starter",
		"HeaderHTML", headerHTML,
		"ContentHTML", contentHTML,
		"FooterHTML", footerHTML,
	)
	if err != nil {
		return "", fmt.Errorf("home page: render page: %w", err)
	}

	return pageHTML, nil
}
