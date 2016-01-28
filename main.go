package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"github.com/apexskier/httpauth"
	"golang.org/x/crypto/bcrypt"
)

var (
	backend     httpauth.GobFileAuthBackend
	aaa         httpauth.Authorizer
	roles       map[string]httpauth.Role
	authFile =  "admin/auth.gob"
	adminLocation = getConfig("AdminLocation")
)

func buildSite() {
	fmt.Println("\n---buildSite---\nCompiling website...")

	out, err := exec.Command("hugo").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\nWebsite compiled!\n---buildSite---\n\n", out)
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

func httpAuth(){
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

func main() {
	initConfig()
	buildSite()
	httpAuth()

	adminLocation := getConfig("AdminLocation")

	http.Handle("/", http.FileServer(http.Dir("./public"))) // serve site output (`/public`) into `/`
	http.Handle(adminLocation+"assets/", http.StripPrefix(adminLocation+"assets/", http.FileServer(http.Dir("./admin/assets"))))
	http.HandleFunc(adminLocation, adminHandler)
	http.HandleFunc(adminLocation+"edit/", editHandler)
	http.HandleFunc(adminLocation+"save/", saveHandler)
	http.HandleFunc(adminLocation+"login-post/", postLogin)
	http.HandleFunc(adminLocation+"login/", loginHandler)
	http.HandleFunc(adminLocation+"logout/", logoutHandler)

	fmt.Println("Starting server on port 3000...")
	fmt.Println("Admin available at:", adminLocation)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
