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

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
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

func dev_iot_device_info(w http.ResponseWriter, r *http.Request) {
	logCurl(w, r)
	response := map[string]interface{}{
		"result": map[string]interface{}{
			"id":            853921,
			"deviceName":    "d_t4_20231221L11054",
			"deviceSecret":  "3799f97b454d27926cc6f3d7cefcd79c",
			"iotInstanceId": "iot-600a5gmp",
			"productKey":    "a54dw3GX0NZ",
			"mqttHost":      "iot-600a5gmp.mqtt.iothub.aliyuncs.com",
			"createdAt":     1706374998078,
			"type":          1,
			"regionId":      "eu-central-1",
		},
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	// Encode and send the response
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func logCurl(w http.ResponseWriter, r *http.Request) {
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
}

func hello(w http.ResponseWriter, r *http.Request) {
	logCurl(w, r)
	client := &http.Client{}
	// Create the response structure
	//response := map[string]interface{}{
	//	"result": map[string]interface{}{
	//		"id":        400019948,
	//		"mac":       "08d1f948e320",
	//		"sn":        "20231221L11054",
	//		"secret":    "d70f3681cb3afb9b",
	//		"timezone":  3.0,
	//		"locale":    "Europe/Moscow",
	//		"shareOpen": 0,
	//		"settings": map[string]interface{}{
	//			"autoWork": 1,
	//		},
	//		"petInTipLimit": 0,
	//	},
	//}

	// Set content type header
	//w.Header().Set("Content-Type", "application/json;charset=utf-8")

	// Encode and send the response
	//err := json.NewEncoder(w).Encode(response)
	//if err != nil {
	//	log.Printf("Error encoding JSON response: %v", err)
	//	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//}

	resp, err := client.Do(r)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()

	log.Println(r.RemoteAddr, " ", resp.Status)

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	//w.Write([]byte("Success"))
}

func main() {
	//http.HandleFunc("/6/t4/dev_iot_device_info", dev_iot_device_info)
	http.HandleFunc("/", hello)
	//http.HandleFunc("/tasks", tasksHandler)
	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	fmt.Printf("Starting server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
