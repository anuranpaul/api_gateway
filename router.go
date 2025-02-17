package main

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"example/API_Gateway/middleware"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	// Middleware (Logging)
	router.Use(chimiddleware.Logger)
	router.Use(middleware.AuthMiddleware)
	router.With(middleware.RequireRole("admin")).Handle("/admin/*", ReverseProxy("http://localhost:5003"))
	router.With(middleware.RequireRole("user")).Handle("/users/*", ReverseProxy("http://localhost:5001"))

	return router
}
