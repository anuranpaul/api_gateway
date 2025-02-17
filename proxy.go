package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"log"
)

// ReverseProxy forwards requests to the backend service
func ReverseProxy(target string) http.Handler {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid backend URL: %s", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = r.URL.Path[len("/users/"):] // Trim prefix before forwarding
		proxy.ServeHTTP(w, r)
	})
}
