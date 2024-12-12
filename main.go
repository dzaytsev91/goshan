package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

var tasks = []Task{
	{
		ID:          "28ba250e-9c11-479f-bcdd-a263cfbd02f0",
		Title:       "First task to do from server",
		IsCompleted: false,
	},
	{
		ID:          "e86444e1-d34e-4ab4-853d-9da9f1d4ae7d",
		Title:       "Second task to do from server",
		IsCompleted: false,
	},
	{
		ID:          "04ae8534-db50-49bb-88e1-90b58620a2eb",
		Title:       "Third task to do from server",
		IsCompleted: false,
	},
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(tasks)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
	return
}

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
	http.HandleFunc("/tasks", tasksHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	fmt.Printf("Starting server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
