package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", getUserHandler)
	mux.HandleFunc("/users/test", getUserHandler)

	server := &http.Server{
		Addr:    ":5001",
		Handler: mux,
	}

	log.Println("User Service running on port 5001")
	log.Fatal(server.ListenAndServe())
}
