package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Admin represents the structure for the admin
type Admin struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Admin functionality handler
func getAdminHandler(w http.ResponseWriter, r *http.Request) {
	admin := Admin{
		ID:   1,
		Name: "Admin User",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admin)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/admin", getAdminHandler)
	mux.HandleFunc("/admin/dashboard", getAdminHandler) // Admin dashboard route

	server := &http.Server{
		Addr:    ":5003", // Admin service runs on port 5003
		Handler: mux,
	}

	log.Println("Admin Service running on port 5003")
	log.Fatal(server.ListenAndServe())
}
