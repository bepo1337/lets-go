package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) initializeRoutes(config *Config) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir(config.staticDir))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fileServer))
	router.GET("/", app.home)
	router.GET("/snippet/view/:id", app.snippetView)
	router.GET("/snippet/create", app.snippetCreateGet)
	router.POST("/snippet/create", app.snippetCreatePost)
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	return standardMiddleware.Then(router)
}
