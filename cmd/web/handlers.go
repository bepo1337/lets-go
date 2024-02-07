package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var HTML_PATH = "./ui/html/"
var HTML_PATH_PAGES = HTML_PATH + "pages/"

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Add("Cache-Control", "max-age=31536000")
	w.Header().Add("Cache-Control", "public")
	//ordering matters here
	templates := []string{
		HTML_PATH + "base.html",
		HTML_PATH + "partials/nav.html",
		HTML_PATH_PAGES + "home.html",
	}
	templateSet, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = templateSet.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "SnippetView func with ID '%d'.", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST") //has to come before WriteHeader
		http.Error(w, "HTTP-Method not allowed!", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(200) //wouldnt need to do this since its default to return 200
	w.Write([]byte("SnippetCreate func2"))
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("SnippetCreate func2"))
}
