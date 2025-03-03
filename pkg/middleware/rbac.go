package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// RequireRole checks if the user has the required role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		
		// Check authentication and role
		isValid, role, _ := CheckAuth(authHeader)  // Add underscore for unused username
		if !isValid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		if role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("Access denied for role: %s", role),
			})
		}

		// Continue to next handler if role matches
		return c.Next()
	}
}
