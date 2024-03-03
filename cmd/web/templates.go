package main

import (
	"github.com/justinas/nosurf"
	"html/template"
	"io/fs"
	"letsgo.bepo1337/internal/models"
	"letsgo.bepo1337/ui"
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
	User            *models.User
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.gohtml",
			"html/partials/*.gohtml",
			page,
		}
		templateSet, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
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
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
