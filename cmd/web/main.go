package main

import (
	"ct-go-web-starter/src"
	"ct-go-web-starter/src/infrastructure/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config.Load()
	src.App()
}
