package main

import (
	// "fmt"
	"github.com/spf13/cast"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/parser"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

type Frontmatter map[string]interface{}

type Content struct {
	Path     string
	Metadata Frontmatter
	Body     string
}

type Dir struct {
	Title string
	Body  []os.FileInfo
}

func (p *Page) save() error {
	filename := p.Title
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func (c *Content) save() error {
	page, err := hugolib.NewPage(c.Path)
	if err != nil {
		return err
	}

	page.SetSourceMetaData(c.Metadata, '+')
	page.SetSourceContent([]byte(c.Body))

	return page.SafeSaveSourceAs(c.Path)
}

func loadPage(title string) (*Page, error) {
	filename := title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func loadContent(filename string) (*Content, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	parser, err := parser.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	rawdata, err := parser.Metadata()
	if err != nil {
		return nil, err
	}

	metadata, err := cast.ToStringMapE(rawdata)
	if err != nil {
		return nil, err
	}

	return &Content{
		Path:     filename,
		Metadata: metadata,
		Body:     string(parser.Content()),
	}, nil
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
	contentDir := getConfig("contentdir")

	if strings.HasPrefix(title, contentDir) {
		c, err := loadContent(title)
		if err != nil {
			c = &Content{Path: title}
		}

		renderTemplateContent(w, "content-edit.html", c)
	} else {
		p, err := loadPage(title)
		if err != nil {
			p = &Page{Title: title}
		}

		renderTemplatePage(w, "edit.html", p)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")

	isAuth(w, r)

	title := r.URL.Path[len(adminLocation+"save/"):]
	body := r.FormValue("body")

	/* if strings.HasPrefix(title, "content") {
		frontmatter := new(Frontmatter)
		for key, value := range r.Form {
			fmt.Println("key:", key)
			fmt.Println("val:", strings.Join(value, ""))
		}
		c := &Content{
			Path:     title,
			Metadata: frontmatter,
			Body:     body,
		}
		c.save()
	} else {*/
	p := &Page{
		Title: title,
		Body:  []byte(body),
	}
	p.save()
	// }
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
