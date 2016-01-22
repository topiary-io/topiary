package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	//	"io/ioutil"
	// "html/template"

	// github.com/spf13/viper for reading config and content front matter
)

func main() {
	fmt.Println("Compiling website...")
	exec.Command("hugo").Output() // rebuild site content
	fmt.Println("Website compiled.")

	http.Handle("/", http.FileServer(http.Dir("./public")))                              // serve site output (`/public`) into `/`
	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("./")))) // serve site input (`/`) into `/admin`

	fmt.Println("Starting server on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
