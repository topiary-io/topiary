package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	// "html/template"
	// github.com/spf13/viper for reading config and content front matter
)

func buildSite() {
	fmt.Println("Compiling website...")
	exec.Command("hugo")
	fmt.Println("Website compiled.")
}

type Page struct {
	Title string
	Body  []byte
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

func adminHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/admin/"):]
	if title == "" {
		title = "./"
	}
	paths, _ := ioutil.ReadDir(title)
	for _, p := range paths {
		if p.IsDir() {
			fmt.Fprintf(w, "<div><a href=\"%s/\">%s</a></div>", p.Name(), p.Name())
		} else {
			fmt.Fprintf(w, "<div>edit: <a href=\"/edit/%s%s\">%s</a></div>", title, p.Name(), p.Name())
		}
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	fmt.Fprintf(w, "<h1>Editing %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\" style=\"width:800px;height:400px;\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		p.Title, p.Title, p.Body)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/admin/", http.StatusFound)
}

func main() {
	buildSite()

	http.Handle("/", http.FileServer(http.Dir("./public"))) // serve site output (`/public`) into `/`
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	fmt.Println("Starting server on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
