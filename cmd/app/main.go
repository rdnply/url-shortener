package main

import (
	"log"

	"github.com/rdnply/url-shortener/internal/app"
)

func main() {
	app, err := app.New(":5000")
	if err != nil {
		log.Fatal(err)
	}

	app.RunServer()
}
