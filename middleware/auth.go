package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	// "github.com/gofiber/fiber/v2"
)

// Secret key for JWT signing
var SecretKey = []byte("your-secret-key")

// ValidateToken parses and validates a JWT token
func ValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("invalid claims")
	}

	return token, claims, nil
}

// CheckAuth validates JWT and extracts the role
func CheckAuth(r *http.Request) (bool, string) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, ""
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, claims, err := ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return false, ""
	}

	// Extract role from claims
	role, ok := claims["role"].(string)
	if !ok {
		return false, ""
	}

	return true, role
}

// AuthMiddleware is the middleware to validate JWT and handle role-based access
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization
		isValid, role := CheckAuth(r)
		if !isValid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add role to context for downstream use
		r = r.WithContext(context.WithValue(r.Context(), "role", role))

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

