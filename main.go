package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	println("got request")
	data, err := json.Marshal(tasks)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
	println("return data %s", data)
	return
}

func hello(w http.ResponseWriter, r *http.Request) {
	// Prepare log data
	fmt.Println("______________________________________________")
	fmt.Printf("%v\n", r)
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("Headers: %s\n", r.Header)
	fmt.Printf("URL: %s\n", r.URL.String())
	switch r.Method {
	case "GET":
		w.Write([]byte("Success"))
	case "POST":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		// Log the raw body
		log.Printf("POST body: %s\n", string(body))
		w.Write([]byte("Success"))
		//if err := r.ParseForm(); err != nil {
		//	fmt.Fprintf(w, "ParseForm() err: %v", err)
		//	return
		//}
		//story := r.FormValue("story")
		//geogess := r.FormValue("geogess")
		//log.Printf("story input: %s, geogess: %s", story, geogess)
		//var story_pass, geogess_pass bool
		//if story == "26270" {
		//	story_pass = true
		//}
		//if strings.ToLower(strings.ReplaceAll(geogess, " ", " ")) == "вандаркхолм" {
		//	geogess_pass = true
		//}
		//
		//if story_pass && geogess_pass {
		//	http.ServeFile(w, r, "image.html")
		//} else {
		//	errorMsg := fmt.Sprintf("First password: %t\nSecond: password: %t\n", story_pass, geogess_pass)
		//	fmt.Printf(errorMsg)
		//	fmt.Fprintf(w, errorMsg)
		//}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
	fmt.Println("______________________________________________")
}

func main() {
	http.HandleFunc("/", hello)
	//http.HandleFunc("/tasks", tasksHandler)
	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	fmt.Printf("Starting server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
