package textinput

import (
	"ct-go-web-starter/internal/web/components/demo"
	"html/template"
)

var Showcase = demo.Page{
	Slug:        "textinput",
	Title:       "Text input",
	Group:       "Inputs",
	Source:      "components/textinput",
	Description: "Single-line text, number, and password inputs, wrapped with a label, hint, and error via the field component.",
	Render:      renderShowcase,
}

// takenUsernames stands in for a database lookup a real server-side check
// would make.
var takenUsernames = map[string]bool{"ada": true}

func renderShowcase(request demo.Request) (template.HTML, error) {
	if request.Query.Has("validate-username") {
		return renderUsernameField(request.BaseHref, request.Query.Get("username"))
	}

	text, err := Render(Options{
		Id:          "showcase-textinput-text",
		Name:        "name",
		Label:       "Name",
		Placeholder: "Ada Lovelace",
		Hint:        "As it should appear on your badge.",
		Required:    true,
	})
	if err != nil {
		return "", err
	}

	number, err := Render(Options{
		Id:    "showcase-textinput-number",
		Name:  "seats",
		Label: "Seats",
		Type:  Number,
		Value: "2",
		Hint:  "How many are in your party?",
	})
	if err != nil {
		return "", err
	}

	invalid, err := Render(Options{
		Id:    "showcase-textinput-invalid",
		Name:  "email",
		Label: "Email",
		Type:  Email,
		Value: "not-an-email",
		Hint:  "We'll only use this to send your receipt.",
	})
	if err != nil {
		return "", err
	}

	username, err := renderUsernameField(request.BaseHref, "")
	if err != nil {
		return "", err
	}

	password, err := Render(Options{
		Id:           "showcase-textinput-password",
		Name:         "password",
		Label:        "Password",
		Type:         Password,
		Autocomplete: "new-password",
		Hint:         "At least 12 characters.",
		Required:     true,
	})
	if err != nil {
		return "", err
	}

	return demo.Frame(demo.FrameOptions{
		Note:    "Leave Name empty then click away, or fix the Email format, to see native validation clear on its own. Username checks against the server on blur (try \"ada\") — no client-side validity state, just a full re-render of the field. Password gets a show/hide toggle.",
		Class:   "max-w-md flex flex-col gap-4",
		Content: text + number + invalid + username + password,
	})
}

// renderUsernameField is used both for the field's initial state and for the
// hx-get response its blur handler requests, so the two never drift apart.
func renderUsernameField(baseHref, value string) (template.HTML, error) {
	var errorText string
	if takenUsernames[value] {
		errorText = "That username is already taken."
	}

	attrs := template.HTMLAttr(
		`hx-get="` + baseHref + `?validate-username" hx-trigger="blur" hx-target="closest .field" hx-swap="outerHTML" hx-sync="closest .field:abort"`,
	)

	return Render(Options{
		Id:    "showcase-textinput-username",
		Name:  "username",
		Label: "Username",
		Type:  Text,
		Value: value,
		Hint:  "Try \"ada\" — it's already taken.",
		Error: errorText,
		Attrs: attrs,
	})
}
