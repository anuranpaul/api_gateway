package main

import (
	"example/API_Gateway/internal/auth"
	"fmt"
)

func main() {
	// Generate tokens for different roles
	adminToken := auth.GenerateJWT("admin_user", "admin")
	userToken := auth.GenerateJWT("regular_user", "user")

	// Print tokens
	fmt.Println("\nAdmin JWT Token:")
	fmt.Println(adminToken)
	fmt.Println("\nUser JWT Token:")
	fmt.Println(userToken)
}
