package main

import (
	"ct-go-web-starter/src"
	"ct-go-web-starter/src/infrastructure/config"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.Load()
	src.App()
}
