package config

import (
	"ct-go-web-starter/src/infrastructure/colorhandler"
	"log/slog"
	"os"
)

func init() {
	// Configure slog with colored output and source location information
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
	handler := colorhandler.New(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
