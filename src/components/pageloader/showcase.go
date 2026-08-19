package pageloader

import (
	"ct-go-web-starter/src/components/demo"
	"html/template"
	"time"
)

var Showcase = demo.Page{
	Slug:        "pageloader",
	Title:       "Page loader",
	Source:      "components/pageloader",
	Description: "A thin progress bar shown at the top of the viewport during a full page navigation.",
	Render:      renderShowcase,
}

// showcaseDelay slows this demo's own response just enough to see the bar
// animate rather than snap straight to complete on a fast local response.
const showcaseDelay = 2500 * time.Millisecond

func renderShowcase(req demo.Request) (template.HTML, error) {
	if req.Query.Get("delay") == "true" {
		time.Sleep(showcaseDelay)
	}

	return demo.Frame(demo.FrameOptions{
		Note:    "Rendered once in the page chrome (see components/page), not instantiated per page. It shows automatically on any plain link click, form submit, or htmx-boosted request, and completes when the navigation lands. This link reloads the current page with an artificial delay so the animation is visible.",
		Content: `<a class="btn btn-primary" href="/showcase/pageloader?delay=true">Navigate to trigger it</a>`,
	})
}
