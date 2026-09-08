package config

import (
	"os"
)

var AppEnv string

func Load() {
	AppEnv = getEnvOr("APP_ENV", "dev")
}

// Port is per-binary: each one listens on PORT, falling back to a different
// default so all three can run side by side locally.
func Port(fallback string) string {
	return getEnvOr("PORT", fallback)
}

func getEnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// MetricsPort is deliberately separate from PORT: /metrics exposes internals and
// must not be published on the port that serves users.
func MetricsPort(fallback string) string {
	return getEnvOr("METRICS_PORT", fallback)
}

// MCPPublicURL is where clients reach the MCP server. The web app cannot derive
// it from its own request — they are separate deployments.
func MCPPublicURL() string {
	return getEnvOr("MCP_PUBLIC_URL", "http://localhost:8082/mcp")
}
