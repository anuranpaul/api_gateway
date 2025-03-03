package test

import (
	"example/API_Gateway/internal/auth"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

// Use the same secret key as in Auth/generate_jwt.go
var secretKey = []byte("your-secret-key")

func TestJWTContents(t *testing.T) {
	// Get token from the auth package
	token := auth.GenerateJWT("testuser", "user")
	
	// Parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	
	if err != nil {
		t.Errorf("Failed to parse token: %v", err)
	}
	
	// Check if token is valid
	if !parsedToken.Valid {
		t.Error("Token is not valid")
	}
	
	// Check claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		if claims["username"] != "testuser" {
			t.Errorf("Expected username 'testuser', got %v", claims["username"])
		}
	} else {
		t.Error("Failed to get claims from token")
	}
}
