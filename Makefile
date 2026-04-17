.PHONY: dev run build

dev:
	air

build:
	go generate ./...
	go build -o build/ .