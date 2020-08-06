package app

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/rdnply/url-shortener/internal/link"
)

func (app *App) loadMainPage(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, app.Templates.main, struct{ NewForm bool }{NewForm: true})
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
}

func (app *App) createLink(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	counter, err := app.CounterStorage.Increment()
	if err != nil {
		app.ServerError(w, err, "")
		return
	}

	encodedLink := app.BaseConvertor.Encode(counter)

	l := &link.Link{URL: url, ShortID: encodedLink, ShortIDInt: counter, Clicks: 0}
	if _, err := app.LinkStorage.AddLink(l); err != nil {
		app.ServerError(w, err, "")
		return
	}

	http.Redirect(w, r, "/stats/"+l.ShortID, http.StatusMovedPermanently)
}

func (app *App) showStats(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "shortID")

	l, err := app.LinkStorage.GetLinkByShortID(shortID)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}

	err = renderTemplate(w, app.Templates.main, struct {
		NewForm bool
		Link    *link.Link
	}{
		NewForm: false,
		Link:    l,
	})
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
}

func (app *App) serverSideRedirect(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "shortID")
	link, err := app.LinkStorage.GetLinkByShortID(shortID)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}

	if _, err := app.LinkStorage.IncrementLinkCounter(link); err != nil {
		app.ServerError(w, err, "")
		return
	}

	http.Redirect(w, r, link.URL, http.StatusSeeOther)
}

func renderTemplate(w io.Writer, tmpl *template.Template, payload interface{}) error {
	err := tmpl.Execute(w, payload)
	if err != nil {
		detail := fmt.Sprintf("can't execute template with name: %v: %v", tmpl.Name(), err)
		return errors.Wrap(err, detail)
	}

	return nil
}
