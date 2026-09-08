package main

import (
	"ct-go-web-starter/internal/infrastructure/compression"
	"ct-go-web-starter/internal/infrastructure/config"
	"ct-go-web-starter/internal/infrastructure/fileserver"
	"ct-go-web-starter/internal/infrastructure/health"
	"ct-go-web-starter/internal/infrastructure/httpserver"
	"ct-go-web-starter/internal/infrastructure/metrics"
	"ct-go-web-starter/internal/infrastructure/reqlog"
	"ct-go-web-starter/internal/service"
	"ct-go-web-starter/internal/web/features/home"
	"ct-go-web-starter/internal/web/features/mcpconnector"
	"ct-go-web-starter/internal/web/features/notfound"
	"ct-go-web-starter/internal/web/features/showcase"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.Load()
	config.InitLogging()

	metrics.Serve(":" + config.MetricsPort("9090"))

	addr := ":" + config.Port("8080")
	handler := reqlog.Middleware()(compression.Middleware()(routes()))

	slog.Info("Web server starting", "addr", "http://localhost"+addr, "version", service.Version)

	if err := httpserver.Run(addr, handler); err != nil {
		slog.Error("Web server error", "error", err)
		os.Exit(1)
	}
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	home.RegisterRoutes(mux)
	showcase.RegisterRoutes(mux)
	mcpconnector.RegisterRoutes(mux)
	health.RegisterRoutes(mux)
	fileserver.RegisterRoutes(mux, "tmp/static/")
	notfound.RegisterRoutes(mux)

	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		reqlog.Skip(r.Context())
		http.Redirect(w, r, "/static/favicon.svg", http.StatusMovedPermanently)
	})

	mux.HandleFunc("/.well-known/appspecific/com.chrome.devtools.json", func(w http.ResponseWriter, r *http.Request) {
		reqlog.Skip(r.Context())
		http.NotFound(w, r)
	})

	return mux
}
