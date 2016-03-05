package main

import (
	//	"fmt"
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

// TODO need a way to lock this to current hugo-site dir
func loadDir(title string) (*Dir, error) {
	dir := title
	body, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return &Dir{Title: title, Body: body}, nil
}

// Might need a better name ; is a location pair where
//	Root = site root (ie /admin/)
//	Dir  = relative file path (ie ./admin)
type Location struct {
	Root string
	Dir  string
}

func serveAdmin(root string, dir string) {
	// TODO some way to standardize roots / dirs (path is wonky)
	// TODO templates should read adminroot (instead of hardcoding /admin/...)
	admin := Location{root, dir}

	// serve admin
	http.HandleFunc(admin.Root, func(w http.ResponseWriter, r *http.Request) {
		adminHandler(w, r, admin)
	})

	// serve assets
	assetroot := admin.Root + "assets/"
	assetdir := admin.Dir + "/assets"
	http.Handle(assetroot, http.StripPrefix(assetroot, http.FileServer(http.Dir(assetdir))))

	// API
	http.HandleFunc(admin.Root+"edit/", func(w http.ResponseWriter, r *http.Request) {
		editHandler(w, r, admin)
	})
	http.HandleFunc(admin.Root+"save/", func(w http.ResponseWriter, r *http.Request) {
		saveHandler(w, r, admin)
	})
	http.HandleFunc(admin.Root+"login/", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(w, r, admin)
	})
	http.HandleFunc(admin.Root+"logout/", func(w http.ResponseWriter, r *http.Request) {
		logoutHandler(w, r, admin)
	})
	http.HandleFunc(admin.Root+"manage-accounts/", func(w http.ResponseWriter, r *http.Request) {
		adminUsersHandler(w, r, admin)
	})
}

func adminHandler(w http.ResponseWriter, r *http.Request, admin Location) {
	isAuth(w, r, admin, "user")

	title := r.URL.Path[len(admin.Root):]
	if title == "" {
		title = "./"
	}

	p, err := loadDir(title)
	if err != nil {
		p = &Dir{Title: title}
	}

	renderTemplateDir(w, "view.html", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, admin Location) {
	// TODO : pass in
	contentDir := getConfig("contentdir")

	isAuth(w, r, admin, "user")

	path := r.URL.Path[len(admin.Root+"edit/"):]

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

func saveHandler(w http.ResponseWriter, r *http.Request, admin Location) {
	contentDir := getConfig("contentdir")

	isAuth(w, r, admin, "user")

	path := r.URL.Path[len(admin.Root+"save/"):]

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
	http.Redirect(w, r, admin.Root, http.StatusFound)
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
