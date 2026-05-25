package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	port := "8080"
	fmt.Printf("Starting local server at http://localhost:%s\n", port)

	// Set up the static file server
	fs := http.FileServer(http.Dir("."))

	// Main multiplexer
	mux := http.NewServeMux()

	// CORS-bypassing proxy handler for MTA feeds
	mux.HandleFunc("/api/mta", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
			return
		}

		// Security: Prevent SSRF by validating that target URL is a valid MTA endpoint
		if !strings.HasPrefix(targetURL, "https://api-endpoint.mta.info/") {
			http.Error(w, "Forbidden target URL: must be an official MTA API endpoint", http.StatusForbidden)
			return
		}

		// Ensure we parse the URL correctly
		parsedURL, err := url.Parse(targetURL)
		if err != nil {
			http.Error(w, "Invalid target URL: "+err.Error(), http.StatusBadRequest)
			return
		}

		req, err := http.NewRequestWithContext(r.Context(), "GET", parsedURL.String(), nil)
		if err != nil {
			http.Error(w, "Failed to create proxy request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Failed to connect to MTA API: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy headers from MTA response
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
	})

	// Static file handler (with correct Content-Type for WASM)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/main.wasm") {
			w.Header().Set("Content-Type", "application/wasm")
		}
		fs.ServeHTTP(w, r)
	})

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
