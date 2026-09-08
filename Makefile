.PHONY: web build up down docker test

ifeq ($(OS),Windows_NT)
  AIR_CONF = .air.windows.toml
else
  AIR_CONF = .air.linux.toml
endif

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS = -s -w -X ct-go-web-starter/internal/service.Version=$(VERSION)

web:
	air -c $(AIR_CONF)

build:
	npm run build-frontend
	go build -ldflags "$(LDFLAGS)" -o build/ ./cmd/...

test:
	go vet ./...
	go test ./...

up: export VERSION := $(VERSION)
up:
	docker compose up --build

down:
	docker compose down

docker:
	docker build -f cmd/web/Dockerfile --build-arg VERSION=$(VERSION) -t ct-go-web-starter-web .
	docker build -f cmd/api/Dockerfile --build-arg VERSION=$(VERSION) -t ct-go-web-starter-api .
	docker build -f cmd/mcp/Dockerfile --build-arg VERSION=$(VERSION) -t ct-go-web-starter-mcp .
	docker image prune -f
