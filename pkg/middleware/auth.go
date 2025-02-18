package middleware

import (
	"fmt"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
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
func CheckAuth(authHeader string) (bool, string) {
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

// AuthMiddleware is the middleware to validate JWT
func AuthMiddleware(c *fiber.Ctx) error {
	// Get authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authorization header",
		})
	}

	// Validate token and get role
	isValid, role := CheckAuth(authHeader)
	if !isValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Add role to context for downstream use
	c.Locals("role", role)

	// Continue to next handler
	return c.Next()
}

