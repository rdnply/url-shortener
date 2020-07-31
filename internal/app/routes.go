package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *App) routes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", app.loadMainPage)
	})

	return r
}
