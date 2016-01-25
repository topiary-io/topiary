package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func buildSite() {
	fmt.Println("\n---buildSite---\nCompiling website...")

	out, err := exec.Command("hugo").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\nWebsite complied!\n---buildSite---\n\n", out)
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

func main() {
	initConfig()
	buildSite()

	adminLocation := getConfig("AdminLocation")

	http.Handle("/", http.FileServer(http.Dir("./public"))) // serve site output (`/public`) into `/`
	http.Handle(adminLocation+"assets/", http.StripPrefix(adminLocation+"assets/", http.FileServer(http.Dir("./admin/assets"))))
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
