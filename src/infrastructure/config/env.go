package config

import (
	"fmt"
	"os"
)

var DirectusURL string

func Load() {
	DirectusURL = mustGetEnv("DIRECTUS_URL")
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("config: %s environment variable is required", key))
	}
	return v
}
