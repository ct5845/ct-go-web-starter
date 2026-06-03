package config

import (
	"fmt"
	"os"
)

var Port string

func Load() {
	Port = getEnvOr("PORT", "8080")
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("config: %s environment variable is required", key))
	}
	return v
}

func getEnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
