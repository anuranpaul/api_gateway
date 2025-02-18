package auth

import (
	"log"
	"time"
	"github.com/dgrijalva/jwt-go"
)

// Secret key for signing the JWT
var SecretKey = []byte("your-secret-key")

// GenerateJWT creates a sample JWT token
func GenerateJWT(username, role string) string {
	// Define claims
	claims := jwt.MapClaims{
		"username": username,
		"role":     role, // ✅ Added role for RBAC
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // Expires in 1 hour
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatalf("Error generating JWT: %s", err)
	}

	return tokenString
}
