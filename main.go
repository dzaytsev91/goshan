package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func hello(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		story := r.FormValue("story")
		geogess := r.FormValue("geogess")
		log.Printf("story input: %s, geogess: %s", story, geogess)
		var story_pass, geogess_pass bool
		if story == "26270" {
			story_pass = true
		}
		if strings.ToLower(strings.ReplaceAll(geogess, " ", " ")) == "вандаркхолм" {
			geogess_pass = true
		}

		if story_pass && geogess_pass {
			http.ServeFile(w, r, "image.html")
		} else {
			errorMsg := fmt.Sprintf("First password: %t\nSecond: password: %t\n", story_pass, geogess_pass)
			fmt.Printf(errorMsg)
			fmt.Fprintf(w, errorMsg)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/", hello)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	fmt.Printf("Starting server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
