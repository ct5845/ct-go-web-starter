package home

import (
	"io"
	"log/slog"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		slog.Warn("Path not found", "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	page, err := renderPage()
	if err != nil {
		slog.Error("Failed to render home page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(page))
}
