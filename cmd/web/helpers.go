package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"
	"runtime/debug"
)

func (app *Application) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *Application) render(w http.ResponseWriter, status int, page string, data *TemplateData) {
	templateSet, ok := app.templates[page]
	if !ok {
		err := fmt.Errorf("the template '%s' doesnt exist", page)
		app.serveError(w, err)
		return
	}
	buffer := new(bytes.Buffer)
	err := templateSet.ExecuteTemplate(buffer, "base", data)
	if err != nil {
		app.serveError(w, err)
		return
	}
	w.WriteHeader(status)
	buffer.WriteTo(w)
}

func (app *Application) decodePostForm(r *http.Request, destination any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(destination, r.PostForm)
	if err != nil {
		if errors.As(err, &form.InvalidDecoderError{}) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *Application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), authenticatedUserId)
}
