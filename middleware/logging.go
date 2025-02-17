package main

import (
	"log"
	"net/http"
	"time"
)

// RequestLogger logs incoming requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("[%s] %s %s - %s", r.Method, r.RequestURI, r.RemoteAddr, duration)
	})
}
