.PHONY: dev build

dev:
	air

build:
	go generate ./...
	go build -o build/ .