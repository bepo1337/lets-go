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

const (
	HTML_PATH           = "./ui/html/"
	HTML_PATH_PAGES     = HTML_PATH + "pages/"
	authenticatedUserId = "authenticateUserID"
)

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

	app.render(w, http.StatusOK, "home.gohtml", templateData)
}

func (app *Application) about(writer http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	app.render(writer, http.StatusOK, "about.gohtml", templateData)

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

	app.render(w, http.StatusOK, "view.gohtml", templateData)

}

func (app *Application) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.gohtml", templateData)
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
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365),
		"expires",
		"Expires not in permitted set")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.gohtml", data)
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

type userRegisterForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *Application) signupUserPost(w http.ResponseWriter, r *http.Request) {
	// Decode
	var userForm userRegisterForm
	err := app.decodePostForm(r, &userForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//Validate
	userForm.CheckField(validator.NotBlank(userForm.Name), "name", "Name cannot be empty!")
	userForm.CheckField(validator.NotBlank(userForm.Email), "email", "E-Mail cannot be empty!")
	userForm.CheckField(validator.NotBlank(userForm.Password), "password", "Password cannot be empty!")
	userForm.CheckField(validator.MinChars(userForm.Password, 2), "password", "Passwords needs to be longer than 1 char!")

	if !userForm.Valid() {
		data := app.newTemplateData(r)
		data.Form = userForm
		app.render(w, http.StatusUnprocessableEntity, "signup.gohtml", data)
		return
	}
	err = app.userModel.Insert(userForm.Name, userForm.Email, userForm.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			userForm.AddFieldError("email", "E-Mail already in use")
			data := app.newTemplateData(r)
			data.Form = userForm
			app.render(w, http.StatusUnprocessableEntity, "signup.gohtml", data)
			return
		} else {
			app.serveError(w, err)
			return
		}
	}
	app.sessionManager.Put(r.Context(), "toast", "User successfully created!")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = &userRegisterForm{}
	app.render(w, http.StatusOK, "signup.gohtml", templateData)
}

type loginUserForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	validator.Validator
}

func (app *Application) loginUserPost(w http.ResponseWriter, r *http.Request) {
	//Decode
	var form = &loginUserForm{}
	err := app.decodePostForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//Validate
	form.Validator.CheckField(validator.NotBlank(form.Email), "email", "E-Mail cant be blank")
	form.Validator.CheckField(validator.NotBlank(form.Password), "password", "Password cant be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.gohtml", data)
		return
	}
	//Authenticate
	id, err := app.userModel.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("E-Mail/Password combination wrong or doesnt exist")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnauthorized, "login.gohtml", data)
		} else {
			app.serveError(w, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serveError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "toast", "Login successful")
	app.sessionManager.Put(r.Context(), authenticatedUserId, id)
	if redirectURLAfterLogin := app.sessionManager.PopString(r.Context(), "redirectURLAfterLogin"); redirectURLAfterLogin != "" {
		http.Redirect(w, r, redirectURLAfterLogin, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = &userRegisterForm{}
	app.render(w, http.StatusOK, "login.gohtml", templateData)
}

func (app *Application) logoutUserPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serveError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), authenticatedUserId)
	app.sessionManager.Put(r.Context(), "toast", "Logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *Application) viewAccount(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.Get(r.Context(), authenticatedUserId).(int)
	userData, err := app.userModel.Get(userId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			app.serveError(w, err)
		}
		return
	}
	templateData := app.newTemplateData(r)
	templateData.User = userData
	app.render(w, http.StatusOK, "account.gohtml", templateData)
}

func (app *Application) changePasswordForm(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = &changePasswordForm{}
	app.render(w, http.StatusOK, "changePassword.gohtml", templateData)
}

func (app *Application) changePasswordPost(w http.ResponseWriter, r *http.Request) {
	//Decode
	var form = &changePasswordForm{}
	err := app.decodePostForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	//validate
	form.Validator.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "Current Password cant be blank")
	form.Validator.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "new password cant be blank")
	form.Validator.CheckField(validator.NotBlank(form.NewPasswordVerification), "newPasswordVerification", "New password verification cant be blank")
	form.Validator.CheckField(validator.Match(form.NewPassword, form.NewPasswordVerification), "newPassword", "Passwords dont match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "changePassword.gohtml", data)
		return
	}
	userId := app.sessionManager.GetInt(r.Context(), authenticatedUserId)
	//check if old pw is correct
	_, err = app.userModel.CorrectPassword(userId, form.CurrentPassword)
	if err != nil {
		if errors.Is(err, models.ErrHashesDontMatch) {
			data := app.newTemplateData(r)
			form.AddFieldError("currentPassword", "Wrong password!")
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "changePassword.gohtml", data)

		} else {
			app.serveError(w, err)
		}
		return
	}
	worked, err := app.userModel.UpdatePassword(userId, form.NewPassword)
	if err != nil {
		app.serveError(w, err)
		return
	}
	if worked {
		app.sessionManager.Put(r.Context(), "toast", "Update password successful")
		http.Redirect(w, r, "/account/view", http.StatusSeeOther)
		return
	} else {
		app.serveError(w, models.ErrNoUpdateFound)
		return
	}
}

type changePasswordForm struct {
	CurrentPassword         string `form:"currentPassword"`
	NewPassword             string `form:"newPassword"`
	NewPasswordVerification string `form:"newPasswordVerification"`
	validator.Validator     `form:"-"`
}
