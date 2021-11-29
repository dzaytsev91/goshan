package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		fmt.Printf("%v\n", r)
		fmt.Printf("tried password: %s\n", r.FormValue("password"))
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		if r.FormValue("password") == "ty pidor" {
			http.ServeFile(w, r, "image.html")
		} else if r.FormValue("password") == "typidor" {
			fmt.Fprintf(w, "It was close!")
		} else {
			fmt.Fprintf(w, "Wrong password!")
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/", hello)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
