package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"log"
)

// ReverseProxy forwards requests to the backend service
func ReverseProxy(target string) http.Handler {
	// Parse the URL of the backend
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid backend URL: %s", err)
	}

	// Create a new reverse proxy that will forward the request to the backend
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't modify the path - let the proxy handle it
		proxy.ServeHTTP(w, r)
	})
}
