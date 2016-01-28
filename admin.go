package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"fmt"
)

type Page struct {
	Title string
	Body  []byte
}

type Dir struct {
	Title string
	Body  []os.FileInfo
}

func (p *Page) save() error {
	filename := p.Title
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func loadDir(title string) (*Dir, error) {
	dir := title
	body, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return &Dir{Title: title, Body: body}, nil
}

func isAuth(w http.ResponseWriter, r *http.Request){
	adminLocation := getConfig("AdminLocation")
	title := r.URL.Path[len(adminLocation):]
	if err := aaa.Authorize(w, r, true); err != nil && title != "login" && title != "login/" {
		fmt.Println(err)
		http.Redirect(w, r, adminLocation+"login/", http.StatusSeeOther)
		return
	}
}

func postLogin(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if err := aaa.Login(w, r, username, password, "/"); err != nil && err.Error() == "already authenticated" {
		http.Redirect(w, r, adminLocation, http.StatusSeeOther)
	} else if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, adminLocation+"/login", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")
	title := r.URL.Path[len(adminLocation+"login/"):]
	p := &Page{Title: title}
	renderTemplatePage(w, "login.html", p)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := aaa.Logout(w, r); err != nil {
		fmt.Println(err)
		// this shouldn't happen
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")

	isAuth(w,r)

	title := r.URL.Path[len(adminLocation):]
	if title == "" {
		title = "./"
	}
	p, err := loadDir(title)
	if err != nil {
		p = &Dir{Title: title}
	}
	/* paths, _ := ioutil.ReadDir(title)
	for _, p := range paths {
		if p.IsDir() {
			fmt.Fprintf(w, "<div><a href=\"%s/\">%s</a></div>", p.Name(), p.Name())
		} else {
			fmt.Fprintf(w, "<div>edit: <a href=\""+adminLocation+"edit/%s%s\">%s</a></div>", title, p.Name(), p.Name())
		}
	} */
	renderTemplateDir(w, "view.html", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")

	isAuth(w,r)

	title := r.URL.Path[len(adminLocation+"edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplatePage(w, "edit.html", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")

	isAuth(w,r)

	title := r.URL.Path[len(adminLocation+"save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	git(title) //this should probably only run if there have been changes
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

func renderTemplateDir(w http.ResponseWriter, tmpl string, p *Dir) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
