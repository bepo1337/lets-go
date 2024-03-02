package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

func secureHeader(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *Application) checkAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.isAuthenticated(r) {
			w.Header().Add("Cache-Control", "no-store")
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serveError(w, fmt.Errorf("%s", err))
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		userSessionId := app.sessionManager.GetInt(r.Context(), authenticatedUserId)
		if userSessionId == 0 {
			next.ServeHTTP(writer, r)
			return
		}
		exists, err := app.userModel.Exists(userSessionId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				app.infoLog.Printf("Didnt find user with sessionId %d in db", userSessionId)
				next.ServeHTTP(writer, r)
				return
			}
			app.serveError(writer, err)
			return
		}
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
			next.ServeHTTP(writer, r)
		}

	})
}
