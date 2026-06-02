package main

import (
	"ct-go-web-starter/src"
	"ct-go-web-starter/src/infrastructure/config"

	"github.com/joho/godotenv"
)

//go:generate npm run build-css
//go:generate go run src/infrastructure/copyassets/copyassets.go

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config.Load()
	src.App()
}
