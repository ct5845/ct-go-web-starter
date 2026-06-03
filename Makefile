.PHONY: web build docker

ifeq ($(OS),Windows_NT)
  AIR_CONF = .air.windows.toml
else
  AIR_CONF = .air.linux.toml
endif

web:
	air -c $(AIR_CONF)

build:
	npm run build-css
	go run ./cmd/copyassets
	go build -o build/ ./cmd/web

docker:
	docker build -f cmd/web/Dockerfile -t ct-go-web-starter .
	docker image prune -f