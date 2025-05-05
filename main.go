package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	targetURL   = "http://api.eu-pet.com" // Your backend server
	proxyPort   = ":80"                   // Proxy listening port
	specialPath = "/t4/dev_device_info"   // Path to intercept
)

func logRequest(r *http.Request) {
	log.Printf(">>> %s %s %s", r.Method, r.URL.Path, r.Proto)
}

func logResponse(resp *http.Response) {
	if resp != nil {
		log.Printf("<<< %s %s", resp.Status, resp.Request.URL.Path)
	}
}

func modifyResponse(resp *http.Response) error {
	// Only modify responses for our special path
	if resp.Request.URL.Path != specialPath {
		return nil
	}

	// Read the original body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse JSON into a generic map
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	// Navigate through the nested structure to find the autowork field
	if result, ok := data["result"].(map[string]interface{}); ok {
		if settings, ok := result["settings"].(map[string]interface{}); ok {
			if autowork, exists := settings["autowork"].(float64); exists {
				log.Printf("Modifying autowork from %.0f to 1", autowork)
				settings["autowork"] = 1
			}
		}
	}

	// Marshal back to JSON
	modifiedBody, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Update response
	resp.Body = io.NopCloser(bytes.NewBuffer(modifiedBody))
	resp.ContentLength = int64(len(modifiedBody))
	resp.Header.Set("Content-Length", string(len(modifiedBody)))

	return nil
}

func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host // Important for virtual hosted sites
		logRequest(req)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		logResponse(resp)
		return modifyResponse(resp)
	}

	return proxy
}

func main() {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal(err)
	}

	proxy := NewReverseProxy(target)
	log.Printf("Starting proxy server on %s", proxyPort)
	log.Fatal(http.ListenAndServe(proxyPort, proxy))
}
