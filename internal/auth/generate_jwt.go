package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key for signing the JWT
var SecretKey = []byte("your-secret-key")

// GenerateJWT creates a JWT token
func GenerateJWT(username, role string) string {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return ""
	}

	return tokenString
}
