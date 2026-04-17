package main

import "ct-go-web-starter/src"

//go:generate npm run build-css
//go:generate go run src/infrastructure/copyassets/copyassets.go

func main() {
	src.App()
}
