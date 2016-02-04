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
	contentDir := getConfig("contentdir")

	isAuth(w, r)

	path := r.URL.Path[len(adminLocation+"edit/"):]

	if strings.HasPrefix(path, contentDir) {

		c, err := new(Page).contentRead(path)

		if err != nil {
			c = &Content{Path: path}
		}

		renderTemplateContent(w, "content-edit.html", c)
	} else {
		p, err := new(Page).Read(path)

		if err != nil {
			p = &Page{Path: path}
		}

		renderTemplatePage(w, "edit.html", p)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")
	contentDir := getConfig("contentdir")

	isAuth(w, r)

	path := r.URL.Path[len(adminLocation+"save/"):]

	if strings.HasPrefix(path, contentDir) {
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

		page.contentUpdate(path, fm, []byte(r.FormValue("body")))
	} else {
		page := &Page{
			Path: path,
			Body: r.FormValue("body"),
		}
		page.Update(path, []byte(r.FormValue("body")))
	}

	git(path) //this should probably only run if there have been changes
	buildSite()
	http.Redirect(w, r, adminLocation, http.StatusFound)
}

var templates = template.Must(template.ParseGlob("admin/templates/*"))

func renderTemplatePage(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
