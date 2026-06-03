package config

import (
	"ct-go-web-starter/src/infrastructure/colorhandler"
	"ct-go-web-starter/src/infrastructure/reqlog"
	"log/slog"
	"os"
)

// InitLogging configures the default slog handler.
// In "prod" it uses structured JSON; otherwise coloured terminal output.
// Must be called after Load() so AppEnv is set.
func InitLogging() {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	var inner slog.Handler
	if AppEnv == "prod" {
		inner = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		inner = colorhandler.New(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(&reqlog.ContextHandler{Inner: inner}))
}
