package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	// "github.com/redis/go-redis/v9"
	"example/API_Gateway/internal/cache"
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

// CheckAuth validates JWT and extracts the role and username
func CheckAuth(authHeader string) (bool, string, string) {
	if authHeader == "" {
		return false, "", ""
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, claims, err := ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return false, "", ""
	}

	// Extract role and username from claims
	role, ok := claims["role"].(string)
	username, ok := claims["username"].(string)
	if !ok {
		return false, "", ""
	}

	return true, role, username
}

// AuthMiddleware is the middleware to validate JWT
func AuthMiddleware(redis *cache.RedisClient) fiber.Handler {

	return func(c *fiber.Ctx) error {
		if c.Path() == "/auth/tokens" {
			return c.Next()
		}
		
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Check Redis cache first
		ctx := context.Background()
		cachedClaims, err := redis.Get(ctx, "token:"+tokenString)
		if err == nil {
			// Token found in cache, set claims
			var claims jwt.MapClaims
			json.Unmarshal([]byte(cachedClaims), &claims)
			c.Locals("role", claims["role"])
			c.Locals("username", claims["username"])
			return c.Next()
		}

		// Validate token and cache if valid
		token, claims, err := ValidateToken(tokenString)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Cache valid token claims
		claimsJSON, _ := json.Marshal(claims)
		redis.Set(ctx, "token:"+tokenString, string(claimsJSON), 24*time.Hour)

		c.Locals("role", claims["role"])
		c.Locals("username", claims["username"])
		return c.Next()
	}
}

