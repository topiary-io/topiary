package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	// "html/template"
)

func buildSite() {
	fmt.Println("\n---buildSite---\nCompiling website...")

	out, err := exec.Command("hugo").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\nWebsite complied!\n---buildSite---\n\n", out)
}

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func git(file string) {
	//need to accommodate for username, email eventually
	cmdName := "git"
	gitAdd := []string{"add", file}
	outadd, err := exec.Command(cmdName, gitAdd...).Output()
	if err != nil {
		fmt.Printf("fatal error on git add\n")
		fmt.Println(err)
	}
	fmt.Printf("\n------git------\ngit add "+file+" %s\n", outadd)

	gitCommit := []string{"commit", "-m", "updating " + file + " from topiary"}
	outci, err := exec.Command(cmdName, gitCommit...).Output()
	if err != nil {
		fmt.Printf("fatal error on git commit\n")
		fmt.Println(err)
	}
	fmt.Printf("git commit "+file+" %s\n------git------\n", outci)
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
	adminLocation := getConfig("AdminLocation")
	title := r.URL.Path[len(adminLocation):]
	if title == "" {
		title = "./"
	}
	paths, _ := ioutil.ReadDir(title)
	for _, p := range paths {
		if p.IsDir() {
			fmt.Fprintf(w, "<div><a href=\"%s/\">%s</a></div>", p.Name(), p.Name())
		} else {
			fmt.Fprintf(w, "<div>edit: <a href=\""+adminLocation+"edit/%s%s\">%s</a></div>", title, p.Name(), p.Name())
		}
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")
	title := r.URL.Path[len(adminLocation+"edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	fmt.Fprintf(w, "<h1>Editing %s</h1>"+
		"<form action=\""+adminLocation+"save/%s\" method=\"POST\">"+
		"<textarea name=\"body\" style=\"width:800px;height:400px;\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>"+
		"<a href=\""+adminLocation+"\">back</a>",
		p.Title, p.Title, p.Body)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")
	title := r.URL.Path[len(adminLocation+"save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	git(title) //this should probably only run if there have been changes
	buildSite()
	http.Redirect(w, r, adminLocation, http.StatusFound)
}

func main() {
	initConfig()
	buildSite()

	adminLocation := getConfig("AdminLocation")

	http.Handle("/", http.FileServer(http.Dir("./public"))) // serve site output (`/public`) into `/`
	http.HandleFunc(adminLocation, adminHandler)
	http.HandleFunc(adminLocation+"edit/", editHandler)
	http.HandleFunc(adminLocation+"save/", saveHandler)

	fmt.Println("Starting server on port 3000...")
	fmt.Println("Admin available at:", adminLocation)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
