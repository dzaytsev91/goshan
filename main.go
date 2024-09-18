package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v\n", r)
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		fmt.Printf("tried password: %s\n", r.FormValue("password"))
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		uzbek := r.FormValue("uzbek")
		devyataev := r.FormValue("devyataev")
		hotel := r.FormValue("hotel")
		var uzbek_pass, devyataev_pass, hotel_pass bool
		if uzbek == "39.834 65.498" {
			uzbek_pass = true
		}
		if devyataev == "54.095 43.266" {
			devyataev_pass = true
		}
		if hotel == "12.958 100.888" {
			hotel_pass = true
		}

		if uzbek_pass && devyataev_pass && hotel_pass {
			http.ServeFile(w, r, "image.html")
		} else {
			errorMsg := fmt.Sprintf("First password: %t\nSecond: password: %t\nThird password: %t", uzbek_pass, devyataev_pass, hotel_pass)
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

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
