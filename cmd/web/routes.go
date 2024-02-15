package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) initializeRoutes(config *Config) http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(config.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	return standardMiddleware.Then(mux)
	//return app.recoverPanic(app.logRequest(secureHeader(mux)))
}
