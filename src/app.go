package src

import (
	"ct-go-web-starter/src/features/home"
	"ct-go-web-starter/src/infrastructure/compression"
	"ct-go-web-starter/src/infrastructure/fileserver"
	_ "embed"
	"log"
	"log/slog"
	"net/http"
)

func App() {
	mux := http.NewServeMux()

	home.RegisterRoutes(mux)
	fileserver.RegisterRoutes(mux, "tmp/static/")

	mux.HandleFunc("/.well-known/appspecific/com.chrome.devtools.json", http.NotFound)

	slog.Info("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", compression.Handler(mux)))
}
