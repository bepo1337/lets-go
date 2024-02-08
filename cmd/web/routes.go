package main

import "net/http"

func (app *Application) initializeRoutes(config *Config) *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(config.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	return mux
}
