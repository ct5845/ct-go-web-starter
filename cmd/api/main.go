package main

import (
	"ct-go-web-starter/internal/api"
	"ct-go-web-starter/internal/infrastructure/compression"
	"ct-go-web-starter/internal/infrastructure/config"
	"ct-go-web-starter/internal/infrastructure/health"
	"ct-go-web-starter/internal/infrastructure/httpserver"
	"ct-go-web-starter/internal/infrastructure/metrics"
	"ct-go-web-starter/internal/infrastructure/reqlog"
	"ct-go-web-starter/internal/service"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.Load()
	config.InitLogging()

	metrics.Serve(":" + config.MetricsPort("9091"))

	addr := ":" + config.Port("8081")
	handler := reqlog.Middleware()(compression.Middleware()(routes()))

	slog.Info("API server starting", "addr", "http://localhost"+addr, "version", service.Version)

	if err := httpserver.Run(addr, handler); err != nil {
		slog.Error("API server error", "error", err)
		os.Exit(1)
	}
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	api.RegisterRoutes(mux)
	health.RegisterRoutes(mux)

	return mux
}
