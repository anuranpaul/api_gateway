package main

import (
	"example/API_Gateway/internal/auth"
	"fmt"
)

func main() {
	// Generate tokens for testing
	adminToken := auth.GenerateJWT("admin_user", "admin")
	userToken := auth.GenerateJWT("testuser", "user")

	// Print tokens
	fmt.Println("\nAdmin Token:")
	fmt.Println(adminToken)
	fmt.Println("\nUser Token:")
	fmt.Println(userToken)
}
