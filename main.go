package main

import (
	"log"
	"net/http"
	// "example/API_Gateway/middleware"
)

func main() {
	router := NewRouter()

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
