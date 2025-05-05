package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

const (
	targetURL    = "http://api.eu-pet.com"
	targetHost   = "api.eu-pet.com"
	proxyPort    = ":8080"
	specialPath  = "/6/t4/dev_device_info"
	specialPath2 = "/6/t4/dev_signup"
)

// logRequest logs the incoming request details
func logRequest(r *http.Request) {
	// Log basic request info
	log.Printf(">>> Request: %s %s %s", r.Method, r.URL.String(), r.Proto)

	// Log headers
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf(">>> Header: %s: %s", name, value)
		}
	}

	// Log body if present
	if r.Body != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf(">>> Error reading request body: %v", err)
		} else {
			log.Printf(">>> Body: %s", string(body))
			// Restore the body so it can be read again
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}
}

// logResponse logs the outgoing response details
func logResponse(resp *http.Response) {
	// Log basic response info
	if resp != nil {
		log.Printf("<<< Response: %s", resp.Status)

		// Log headers
		for name, values := range resp.Header {
			for _, value := range values {
				log.Printf("<<< Header: %s: %s", name, value)
			}
		}

		// Log body if present
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("<<< Error reading response body: %v", err)
		} else {
			log.Printf("<<< Body: %s", string(body))
			// Restore the body so it can be read again
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
		}
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
	resp.Body.Close()

	// Check if response is JSON before parsing
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		log.Printf("Response is not JSON (Content-Type: %s), skipping modification", contentType)
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		resp.ContentLength = int64(len(body))
		resp.Header.Set("Content-Length", string(len(body)))
		return nil
	}

	// Parse JSON into a generic map
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("JSON parse error: %v", err)
		// Return original response if JSON is invalid
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		resp.ContentLength = int64(len(body))
		resp.Header.Set("Content-Length", string(len(body)))
		return nil
	}

	// Navigate through the nested structure to find the autowork field
	if result, ok := data["result"].(map[string]interface{}); ok {
		if settings, ok := result["settings"].(map[string]interface{}); ok {
			if autowork, exists := settings["autoWork"].(float64); exists {
				log.Printf("Modifying autowork from %.0f to 1", autowork)
				settings["autoWork"] = 1
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
	resp.Header.Set("Content-Length", strconv.Itoa(len(modifiedBody)))
	logResponse(resp)

	return nil
}

func NewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
		logRequest(req)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		logResponse(resp)
		return modifyResponse(resp)
	}

	// Add error handling for proxy errors
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		rw.WriteHeader(http.StatusBadGateway)
		_, _ = rw.Write([]byte("Proxy error: " + err.Error()))
	}

	return proxy
}

func hostSpecificHandler(proxy http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Host != targetHost {
			log.Println("Ban request with wrong domain, %s", r.URL.Host)
			return
		}
		// Only proxy requests for the specific host
		log.Printf("Proxying request for host: %s", r.Host)
		proxy.ServeHTTP(w, r)
		return
	})
}

func main() {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal(err)
	}

	proxy := NewReverseProxy(target)
	handler := hostSpecificHandler(proxy)
	log.Printf("Starting proxy server on %s", proxyPort)
	log.Fatal(http.ListenAndServe(proxyPort, handler))
}
