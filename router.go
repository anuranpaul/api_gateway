package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	// Middleware (Logging)
	router.Use(RequestLogger)

	// Define routes
	router.Handle("/users/*", ReverseProxy("http://localhost:5001"))
	router.Handle("/orders/*", ReverseProxy("http://localhost:5002"))

	return router
}
