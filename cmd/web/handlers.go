package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"letsgo.bepo1337/internal/models"
	"letsgo.bepo1337/internal/validator"
	"net/http"
	"strconv"
)

const HTML_PATH = "./ui/html/"
const HTML_PATH_PAGES = HTML_PATH + "pages/"

var permittedExpireValues = [3]int{1, 7, 365}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
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
	params := httprouter.ParamsFromContext(r.Context())
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

func (app *Application) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.html", templateData)
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "Title cant be blank")
	form.CheckField(validator.WithinMaxChars(form.Title, 100),
		"title",
		"Title cant be greater than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "Content cant be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365),
		"expires",
		"Expires not in permitted set")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}
	id, err := app.snippetModel.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serveError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "toast", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` //decoder ignores this field
}
