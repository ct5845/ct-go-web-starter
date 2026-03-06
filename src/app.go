package src

import (
	"ct-go-web-starter/src/features/home"
	"ct-go-web-starter/src/infrastructure/fileserver"
	"log"
	"log/slog"
	"net/http"
)

func App() {
	// Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", home.Handler)

	// Handle Chrome DevTools well-known endpoint to keep logs clean
	mux.HandleFunc("/.well-known/appspecific/com.chrome.devtools.json", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Static files with ETag caching
	cachedFS := fileserver.NewCachedFileServer("tmp/static/")
	mux.Handle("/static/", http.StripPrefix("/static/", cachedFS))

	slog.Info("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}
