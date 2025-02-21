package routes

import (
	"example/API_Gateway/internal/handlers"
	"example/API_Gateway/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userHandler *handlers.UserHandler) {
	users := app.Group("/api/users")
	
	// Apply auth middleware to the group
	users.Use(func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		isValid, role, username := middleware.CheckAuth(authHeader)
		if !isValid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("role", role)
		c.Locals("username", username)
		return c.Next()
	})

	// Routes
	users.Get("/", middleware.RequireRole("admin"), userHandler.GetAllUsers)
	users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", middleware.RequireRole("admin"), userHandler.UpdateUser)
	users.Patch("/:id", userHandler.PatchUser)
	users.Delete("/:id", userHandler.DeleteUser)
} 