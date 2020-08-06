package main

import (
	"io"
	"log"

	"github.com/rdnply/url-shortener/internal/app"
)

func main() {
	app, closers, err := app.New(":5000")
	if err != nil {
		log.Fatal(err)
	}
	defer handleClosers(closers)

	app.RunServer()
}

func handleClosers(m map[string]io.Closer) {
	for n, c := range m {
		if err := c.Close(); err != nil {
			log.Printf("can't close %q: %s", n, err)
		}
	}
}
