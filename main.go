package main

import (
	// "fmt"
	"log"
	"net/http"
	"os/exec"
	//	"io/ioutil"
)

func main() {
	exec.Command("hugo").Output()                                                        // rebuild site content
	http.Handle("/", http.FileServer(http.Dir("./public")))                              // serve site output (`/public`) into `/`
	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("./")))) // serve site input (`/`) into `/admin`
	err := http.ListenAndServe(":3000", nil)                                             // start server on port 3000
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
