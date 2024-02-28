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

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.signupUserPost))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.signupUserForm))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.loginUserPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.loginUserForm))

	protected := dynamic.Append(app.checkAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreateGet))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.logoutUserPost))

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	return standardMiddleware.Then(router)
}
