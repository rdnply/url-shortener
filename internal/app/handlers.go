package app

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func (app *App) loadMainPage(w http.ResponseWriter, r *http.Request) {
	err := renderTemplate(w, app.Templates.main, struct{ NewForm bool }{NewForm: true})
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
}

func (app *App) createLink(w http.ResponseWriter, r *http.Request) {

}

func renderTemplate(w io.Writer, tmpl *template.Template, payload interface{}) error {
	err := tmpl.Execute(w, payload)
	if err != nil {
		detail := fmt.Sprintf("can't execute template with name: %v: %v", tmpl.Name(), err)
		return errors.Wrap(err, detail)
	}

	return nil
}
