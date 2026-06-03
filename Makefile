.PHONY: dev run build

dev:
	air

build:
	npm run build-css
	go run ./cmd/copyassets
	go build -o build/ ./cmd/web