package mcpconnector

import (
	"ct-go-web-starter/internal/infrastructure/config"
	"ct-go-web-starter/internal/infrastructure/reqlog"
	"ct-go-web-starter/internal/mcpserver"
	"ct-go-web-starter/internal/web/components/component"
	"ct-go-web-starter/internal/web/components/layoutswitch"
	"ct-go-web-starter/internal/web/components/page"
	"ct-go-web-starter/internal/web/features/nav"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

//go:embed mcpconnector.html
var mcpconnectorHTML string
var mcpconnectorTmpl = component.New("mcpconnector.html", mcpconnectorHTML)

const (
	title       = "MCP Connector"
	description = "Connect an LLM client to this application's tools."
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /connector", HandleGet)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	defer reqlog.Track(r.Context(), "mcpconnector.HandleGet", "")()

	rendered, err := render()
	if err != nil {
		slog.Error("Failed to render connector page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(rendered))
}

func render() (template.HTML, error) {
	remoteURL := config.MCPPublicURL()

	content, err := mcpconnectorTmpl.Render(struct {
		Title        string
		Description  string
		ServerName   string
		RemoteURL    string
		RemoteConfig string
	}{
		Title:        title,
		Description:  description,
		ServerName:   mcpserver.ServerName,
		RemoteURL:    remoteURL,
		RemoteConfig: remoteConfig(remoteURL),
	})
	if err != nil {
		return "", fmt.Errorf("connector page: render content: %w", err)
	}

	navigation, err := nav.Render("connector")
	if err != nil {
		return "", fmt.Errorf("connector page: render navigation: %w", err)
	}

	return layoutswitch.RenderPage(page.Options{
		Title:           title,
		MetaDescription: description,
	}, layoutswitch.Options{
		Content:    content,
		BottomTabs: navigation.Footer,
		SideBar:    navigation.SideBar,
	})
}

func remoteConfig(url string) string {
	connector := map[string]any{
		"mcpServers": map[string]any{
			mcpserver.ServerName: map[string]string{
				"type": "http",
				"url":  url,
			},
		},
	}

	encoded, err := json.MarshalIndent(connector, "", "  ")
	if err != nil {
		return ""
	}

	return string(encoded)
}
