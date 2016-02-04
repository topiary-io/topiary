package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Dir struct {
	Title string
	Body  []os.FileInfo
}

func loadDir(title string) (*Dir, error) {
	dir := title
	body, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return &Dir{Title: title, Body: body}, nil
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")

	isAuth(w, r)

	title := r.URL.Path[len(adminLocation):]
	if title == "" {
		title = "./"
	}
	p, err := loadDir(title)
	if err != nil {
		p = &Dir{Title: title}
	}

	renderTemplateDir(w, "view.html", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")

	isAuth(w, r)

	title := r.URL.Path[len(adminLocation+"edit/"):]

	c, err := new(Page).Read(title)

	if err != nil {
		c = &Content{Path: title}
	}

	renderTemplateContent(w, "content-edit.html", c)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")

	isAuth(w, r)

	path := r.URL.Path[len(adminLocation+"save/"):]
	var page Page

	r.ParseForm()

	fm := make(map[string]interface{})

	for key, values := range r.Form {
		for _, value := range values {
			if strings.HasPrefix(key, "fm.") {
				fm[key[3:]] = value
			}
		}
	}

	page.Update(path, fm, []byte(r.FormValue("body")))

	git(path) //this should probably only run if there have been changes
	buildSite()
	http.Redirect(w, r, adminLocation, http.StatusFound)
}

var templates = template.Must(template.ParseGlob("admin/templates/*"))

func renderTemplateContent(w http.ResponseWriter, tmpl string, c *Content) {
	err := templates.ExecuteTemplate(w, tmpl, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateDir(w http.ResponseWriter, tmpl string, p *Dir) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
