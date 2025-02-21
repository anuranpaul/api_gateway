package handlers

import (
	"context"
	"example/API_Gateway/internal/auth"
	"example/API_Gateway/internal/cache"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	redis *cache.RedisClient
}

func NewAuthHandler(redis *cache.RedisClient) *AuthHandler {
	return &AuthHandler{redis: redis}
}

func (h *AuthHandler) GetTokens(c *fiber.Ctx) error {
	adminToken := auth.GenerateJWT("admin_user", "admin")
	userToken := auth.GenerateJWT("testuser", "user")

	// Cache tokens
	ctx := context.Background()
	err := h.redis.Set(ctx, "token:admin_user", adminToken, 24*time.Hour)
	if err != nil {
		fmt.Printf("Error caching admin token: %v\n", err)
	}
	err = h.redis.Set(ctx, "token:testuser", userToken, 24*time.Hour)
	if err != nil {
		fmt.Printf("Error caching user token: %v\n", err)
	}

	return c.JSON(fiber.Map{
		"adminToken": adminToken,
		"userToken":  userToken,
	})
} 