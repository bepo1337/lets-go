package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"letsgo.bepo1337/ui"
	"net/http"
)

func (app *Application) initializeRoutes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.signupUserPost))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.signupUserForm))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.loginUserPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.loginUserForm))

	protected := dynamic.Append(app.checkAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreateGet))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.logoutUserPost))
	router.Handler(http.MethodGet, "/account/view", protected.ThenFunc(app.viewAccount))

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)
	return standardMiddleware.Then(router)
}
