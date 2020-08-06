package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *App) routes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", app.loadMainPage)
		r.Post("/new", app.createLink)
		r.Get("/s/{shortID}", app.serverSideRedirect)
	})

	return r
}
