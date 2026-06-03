package config

import (
	"fmt"
	"os"
)

// example: var MyKey string

func Load() {
	// example: MyKey = mustGetEnv("MY_KEY")
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("config: %s environment variable is required", key))
	}
	return v
}
