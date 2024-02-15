package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"letsgo.bepo1337/internal/models"
	"net/http"
	"strconv"
)

const HTML_PATH = "./ui/html/"
const HTML_PATH_PAGES = HTML_PATH + "pages/"

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	w.Header().Add("Cache-Control", "max-age=31536000")
	w.Header().Add("Cache-Control", "public")
	snippets, err := app.snippetModel.LatestTen()
	if err != nil {
		app.serveError(w, err)
		return
	}
	templateData := app.newTemplateData(r)
	templateData.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", templateData)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippetModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serveError(w, err)
		}
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Snippet = snippet
	app.render(w, http.StatusOK, "view.html", templateData)

}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST") //has to come before WriteHeader
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		app.serveError(w, err)
		return
	}
	var result map[string]string
	json.Unmarshal(bodyBytes, &result)
	expiresAsInt, err := strconv.Atoi(result["expires"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	id, err := app.snippetModel.Insert(result["title"], result["content"], expiresAsInt)
	if err != nil {
		app.serveError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
