package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"letsgo.bepo1337/internal/models"
	"net/http"
	"strconv"
)

const HTML_PATH = "./ui/html/"
const HTML_PATH_PAGES = HTML_PATH + "pages/"

func (app *Application) home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
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

func (app *Application) snippetCreateGet(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Write([]byte("Form data here..."))
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
