package main

import (
	"html/template"
	"letsgo.bepo1337/internal/models"
	"net/http"
	"path/filepath"
	"time"
)

type TemplateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
	Form        any
	Toast       string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(HTML_PATH + "pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).Funcs(functions).ParseFiles(HTML_PATH + "base.html")
		if err != nil {
			return nil, err
		}
		templateSet, err = templateSet.ParseGlob(HTML_PATH + "partials/*.html")

		templateSet, err = templateSet.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = templateSet
	}
	return cache, nil
}

func (app *Application) newTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{CurrentYear: time.Now().Year()}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
