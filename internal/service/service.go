// Package service holds the operations every transport exposes.
//
// It must not import anything else under internal/ — the web, api and mcp
// packages all point inward at this one. If service ever needs to import a
// transport or a store, the dependency has gone the wrong way round.
package service

import (
	"errors"
	"strings"
)

// Version is stamped at build time with
// -ldflags "-X ct-go-web-starter/internal/service.Version=<version>".
var Version = "dev"

var ErrEmptyName = errors.New("name must not be empty")

type Health struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func CheckHealth() Health {
	return Health{Status: "ok", Version: Version}
}

type Greeting struct {
	Message string `json:"message"`
}

func Greet(name string) (Greeting, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Greeting{}, ErrEmptyName
	}
	return Greeting{Message: "Hello, " + name + "!"}, nil
}
