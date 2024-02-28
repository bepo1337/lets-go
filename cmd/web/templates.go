package main

import (
	"github.com/justinas/nosurf"
	"html/template"
	"letsgo.bepo1337/internal/models"
	"net/http"
	"path/filepath"
	"time"
)

type TemplateData struct {
	IsAuthenticated bool
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	CurrentYear     int
	Form            any
	Toast           string
	CSRFToken       string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(HTML_PATH + "pages/*.gohtml")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(HTML_PATH + "base.gohtml")
		if err != nil {
			return nil, err
		}
		templateSet, err = templateSet.ParseGlob(HTML_PATH + "partials/*.gohtml")

		templateSet, err = templateSet.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = templateSet
	}
	return cache, nil
}

func (app *Application) newTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		IsAuthenticated: app.isAuthenticated(r),
		CurrentYear:     time.Now().Year(),
		Toast:           app.sessionManager.PopString(r.Context(), "toast"),
		CSRFToken:       nosurf.Token(r),
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
