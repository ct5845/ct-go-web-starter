package main

import (
	"context"
	"ct-go-web-starter/src/features/home"
	"ct-go-web-starter/src/infrastructure/compression"
	"ct-go-web-starter/src/infrastructure/config"
	"ct-go-web-starter/src/infrastructure/fileserver"
	"ct-go-web-starter/src/infrastructure/reqlog"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.Load()
	config.InitLogging()
	run()
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	home.RegisterRoutes(mux)
	fileserver.RegisterRoutes(mux, "tmp/static/")

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

func run() {
	mux := routes()

	handler := reqlog.Middleware()(compression.Middleware()(mux))

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler,
	}

	slog.Info("Server starting", "addr", "http://localhost:"+config.Port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")
}
