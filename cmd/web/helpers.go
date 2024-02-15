package main

import (
	"bytes"
	"fmt"
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
