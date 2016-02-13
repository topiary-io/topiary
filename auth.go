package main

import (
	"fmt"
	"github.com/apexskier/httpauth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"html/template"
)

var (
	backend  httpauth.GobFileAuthBackend
	aaa      httpauth.Authorizer
	roles    map[string]httpauth.Role
	authFile = "admin/auth.gob"
)

func initAuth() {
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

func isAuth(w http.ResponseWriter, r *http.Request, role string) {
	// first check if you are logged in to the admin if not then redirect to login page
	adminLocation := getConfig("AdminLocation")
	title := r.URL.Path[len(adminLocation):]
	if err := aaa.Authorize(w, r, true); err != nil && title != "login/" {
		fmt.Println(err)
		http.Redirect(w, r, adminLocation+"login/", http.StatusSeeOther)
		return
	}
	// next check if you have the required role
	if err_role := aaa.AuthorizeRole(w,r,role,false); err_role != nil {
		fmt.Println(err_role)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	adminLocation := getConfig("AdminLocation")
	if r.Method == "POST" {
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if err := aaa.Login(w, r, username, password, "/"); err != nil && err.Error() == "already authenticated" {
			http.Redirect(w, r, adminLocation, http.StatusSeeOther)
		} else if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, adminLocation+"/login/", http.StatusSeeOther)
		}
	} else {
		title := r.URL.Path[len(adminLocation+"login/"):]
		p := &Content{Path: title}
		renderTemplateContent(w, "login.html", p)
	}
}

func adminUsersHandler(w http.ResponseWriter, r *http.Request) {
	isAuth(w,r,"admin")
	if r.Method == "POST" {
		var user httpauth.UserData
		user.Username = r.PostFormValue("username")
		user.Email = r.PostFormValue("email")
		password := r.PostFormValue("password")
		user.Role = r.PostFormValue("role")
		if err := aaa.Register(w, r, user, password); err != nil {
		// maybe something
		}
	}

	if user, err := aaa.CurrentUser(w,r); err == nil {
		type data struct {
			User  httpauth.UserData
			Roles map[string]httpauth.Role
			Users []httpauth.UserData
			Msg   []string
		}
		messages := aaa.Messages(w, r)
		users, err := backend.Users()
		if err != nil {
			panic(err)
		}
		d := data{User: user, Roles: roles, Users: users, Msg: messages}
		var templates = template.Must(template.ParseGlob("admin/templates/*"))
		t_err := templates.ExecuteTemplate(w, "manage-accounts.html", d)
		if t_err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
