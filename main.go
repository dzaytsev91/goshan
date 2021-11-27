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
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		if r.FormValue("name") == "ty pidor" {
		    w.Header().Set("Content-Type", "text/html; charset=utf-8")
            fmt.Fprintf(w, "<img src='assets/pidrila.jpg' width='500' alt=''><br><a href='https://clck.ru/Z4pzH'>Лучше не кликай</a>")
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
