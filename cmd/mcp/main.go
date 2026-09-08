package main

import (
	"ct-go-web-starter/internal/infrastructure/config"
	"ct-go-web-starter/internal/infrastructure/health"
	"ct-go-web-starter/internal/infrastructure/httpserver"
	"ct-go-web-starter/internal/infrastructure/metrics"
	"ct-go-web-starter/internal/infrastructure/reqlog"
	"ct-go-web-starter/internal/mcpserver"
	"ct-go-web-starter/internal/service"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	godotenv.Load()
	config.Load()
	config.InitLogging()

	metrics.Serve(":" + config.MetricsPort("9092"))

	addr := ":" + config.Port("8082")

	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return mcpserver.New() },
		&mcp.StreamableHTTPOptions{Stateless: true},
	)

	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)
	health.RegisterRoutes(mux)

	slog.Info("MCP server starting", "addr", "http://localhost"+addr+"/mcp", "version", service.Version)

	if err := httpserver.Run(addr, reqlog.Middleware()(mux)); err != nil {
		slog.Error("MCP server error", "error", err)
		os.Exit(1)
	}
}
