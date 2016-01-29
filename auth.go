package main

import (
	"fmt"
	"net/http"
	"github.com/apexskier/httpauth"
	"golang.org/x/crypto/bcrypt"
)

var (
	backend     httpauth.GobFileAuthBackend
	aaa         httpauth.Authorizer
	roles       map[string]httpauth.Role
	authFile =  "admin/auth.gob"
)

func initAuth(){
	var err error

	// authFile must exist in site home.
	// could intruduce a way to dynamically create one or just default install it
	backend, err = httpauth.NewGobFileAuthBackend(authFile)
	if err != nil {
		panic(err)
	}

	// create some default roles
	roles = make(map[string]httpauth.Role)
	roles["user"] = 30
	roles["admin"] = 80
	aaa, err = httpauth.NewAuthorizer(backend, []byte("cookie-encryption-key"), "user", roles)

	// create a default user
	hash, err := bcrypt.GenerateFromPassword([]byte("adminadmin"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	defaultUser := httpauth.UserData{Username: "admin", Email: "admin@localhost", Hash: hash, Role: "admin"}
	err = backend.SaveUser(defaultUser)
	if err != nil {
		panic(err)
	}
}
 
func isAuth(w http.ResponseWriter, r *http.Request){
	adminLocation := getConfig("AdminLocation")
	title := r.URL.Path[len(adminLocation):]
	if err := aaa.Authorize(w, r, true); err != nil && title != "login/" {
		fmt.Println(err)
		http.Redirect(w, r, adminLocation+"login/", http.StatusSeeOther)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		adminLocation := getConfig("AdminLocation")
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if err := aaa.Login(w, r, username, password, "/"); err != nil && err.Error() == "already authenticated" {
			http.Redirect(w, r, adminLocation, http.StatusSeeOther)
		} else if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, adminLocation+"/login/", http.StatusSeeOther)
		}
	} else {
		adminLocation := getConfig("AdminLocation")
		title := r.URL.Path[len(adminLocation+"login/"):]
		p := &Page{Title: title}
		renderTemplatePage(w, "login.html", p)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")
	if err := aaa.Logout(w, r); err != nil {
		fmt.Println(err)
		// this shouldn't happen
		return
	}
	http.Redirect(w, r, adminLocation+"/login/", http.StatusSeeOther)
}

