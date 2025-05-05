package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	// Create a buffer to build our cURL command
	var curlCmd bytes.Buffer

	// Start with the basic curl command and method
	curlCmd.WriteString("curl -X ")
	curlCmd.WriteString(r.Method)
	curlCmd.WriteString(" \\\n")

	// Add the URL
	curlCmd.WriteString("  '")
	curlCmd.WriteString("http://") // or https:// based on your server
	curlCmd.WriteString(r.Host)
	curlCmd.WriteString(r.URL.String())
	curlCmd.WriteString("' \\\n")

	// Add all headers
	for name, values := range r.Header {
		for _, value := range values {
			curlCmd.WriteString("  -H '")
			curlCmd.WriteString(name)
			curlCmd.WriteString(": ")
			curlCmd.WriteString(value)
			curlCmd.WriteString("' \\\n")
		}
	}

	// Handle request body for POST, PUT, PATCH
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		// Read the body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
		} else {
			// Restore the body
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(bodyBytes) > 0 {
				// Escape single quotes in the body
				bodyStr := strings.ReplaceAll(string(bodyBytes), "'", `'\''`)
				curlCmd.WriteString("  --data-raw '")
				curlCmd.WriteString(bodyStr)
				curlCmd.WriteString("' \\\n")
			}
		}
	}

	// Remove the trailing backslash and newline
	curlString := strings.TrimSuffix(curlCmd.String(), " \\\n")

	// Log the complete cURL command
	log.Println("Incoming request as cURL:")
	log.Println(curlString)

	// Your normal response
	w.Write([]byte("Success"))
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
