package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
)

// TODO - should gzip files
func mithrilAdminLocation() string {
	return "/mithril-ui/"
}

// structs
type Config struct {
	AdminLocation string `json:"adminLocation"`
}

type Link struct {
	Text string `json:"text"`
	URI  string `json:"uri"`
}

type LinkGroup struct {
	Name    string `json:"name"`
	Options []Link `json:"options"`
}

type Document struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type APIResponse struct {
	Message string `json:"message"`
}

// methods
func (d *Document) save() error {
	fn := d.Filename
	return ioutil.WriteFile(fn, []byte(d.Content), 0600)
}

// Assumes JSON transfer methods
func setHeaders(w http.ResponseWriter, r *http.Request) {
	// TODO move into mux
	w.Header().Set("Access-Control-Allow-Origin", "*") // should be from config
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Methods", "POST")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "OPTIONS" {
		return
	}
}

// almost all api calls should send a json response
func sendJSONResponse(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func readDocument(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, r)

	adminLocation := mithrilAdminLocation()
	fn := r.URL.Path[len(adminLocation+"api/read/"):]

	body, err := ioutil.ReadFile(fn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	doc := Document{
		fn, string(body),
	}

	sendJSONResponse(w, doc)
}

func writeDocument(w http.ResponseWriter, r *http.Request) {
	setHeaders(w, r)

	var d Document
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		fmt.Println(err)
	}

	if d.Filename == "" {
		msg := APIResponse{"no filename"}
		sendJSONResponse(w, msg)
		return
	}

	if d.Content == "" {
		msg := APIResponse{"no file content"}
		sendJSONResponse(w, msg)
		return
	}

	d.save()
	msg := APIResponse{d.Filename + " saved"}
	sendJSONResponse(w, msg)
}

func serveMithrilAdmin() {
	adminLocation := mithrilAdminLocation()

	tmplBytes, err := Asset("admin/ui/data/index.html")
	if err != nil {
		fmt.Println("Template error:", err)
	}

	tmpl := template.Must(template.New("index").Parse(string(tmplBytes)))
	cfg := Config{adminLocation}

	//	fs := http.FileServer(http.Dir("." + getConfig("AdminLocation") + "ui"))
	//	http.Handle(adminLocation+"ui/", http.StripPrefix(adminLocation+"ui", fs))

	http.HandleFunc(path.Join(adminLocation, "ui/bundle.js"), func(w http.ResponseWriter, r *http.Request) {
		data, err := Asset("admin/ui/data/bundle.js")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(data)
	})

	http.HandleFunc(path.Join(adminLocation, "ui/main.css"), func(w http.ResponseWriter, r *http.Request) {
		data, err := Asset("admin/ui/data/main.css")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text/css")
		w.Write(data)
	})

	http.HandleFunc(adminLocation, func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// should be constructed from config file (json/toml/?)?
	// have a function which builds an "initialize" data endpoint?
	sideBar := []LinkGroup{
		LinkGroup{
			"Content",
			[]Link{
				Link{"Posts", adminLocation + "content/posts"},
				Link{"Pages", adminLocation + "content/pages"},
				Link{"Media", adminLocation + "static/media"},
			},
		},
		LinkGroup{
			"Theme",
			[]Link{
				Link{"Layouts", adminLocation + "layouts"},
				Link{"Assets", adminLocation + "static"},
				Link{"Data", adminLocation + "data"},
			},
		},
		LinkGroup{
			"Settings",
			[]Link{
				Link{"Configuration", adminLocation + "edit/config.toml"},
				Link{"Archetypes", adminLocation + "archetypes"},
				Link{"Users", adminLocation + "users"},
			},
		},
		LinkGroup{
			"Test",
			[]Link{
				Link{"Loader", adminLocation + "loader"},
			},
		},
	}

	http.HandleFunc(adminLocation+"api/side-bar", func(w http.ResponseWriter, r *http.Request) {
		setHeaders(w, r)
		sendJSONResponse(w, sideBar)
	})

	http.HandleFunc(adminLocation+"api/read/", readDocument)
	http.HandleFunc(adminLocation+"api/write/", writeDocument)

	// will send a 501 Not Implemented page
	http.HandleFunc(adminLocation+"api/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "501")
	})

	fmt.Println("Admin(mithril-ui) Location:", adminLocation)
}
