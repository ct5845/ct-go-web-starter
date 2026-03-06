package main

import (
	"ct-go-web-starter/src/shared/utils"
	"log/slog"
	"os"
)

func main() {
	if err := os.MkdirAll("tmp/static", 0755); err != nil {
		slog.Error("Failed to create tmp/static directory", "error", err)
		panic(err)
	}

	if err := copyAssets(); err != nil {
		panic(err)
	}

	if err := copyAlpineJS(); err != nil {
		panic(err)
	}

	if err := copyHTMX(); err != nil {
		panic(err)
	}
}

func copyAssets() error {
	err := utils.CopyDir("src/static", "tmp/static")
	if err != nil {
		slog.Error("Failed to copy assets", "error", err)
		return err
	}

	slog.Info("Copied assets to tmp/static")
	return nil
}

func copyAlpineJS() error {
	err := utils.CopyFile("node_modules/alpinejs/dist/cdn.min.js", "tmp/static/alpine.min.js")
	if err != nil {
		slog.Error("Failed to copy Alpine.js", "error", err)
		return err
	}

	slog.Info("Copied Alpine.js to tmp/static")
	return nil
}

func copyHTMX() error {
	err := utils.CopyFile("node_modules/htmx.org/dist/htmx.min.js", "tmp/static/htmx.min.js")
	if err != nil {
		slog.Error("Failed to copy HTMX", "error", err)
		return err
	}

	slog.Info("Copied HTMX to tmp/static")
	return nil
}
