package routes

import (
	"example/API_Gateway/internal/handlers"
	"example/API_Gateway/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userHandler *handlers.UserHandler) {
	// User management routes
	users := app.Group("/api/users")
	
	// Protected routes (require authentication)
	users.Use(middleware.AuthMiddleware)

	// Routes accessible by all authenticated users
	users.Post("/", userHandler.CreateUser)           // Create user (role restriction in handler)
	users.Get("/:id", userHandler.GetUser)           // Get user (self or admin only)
	users.Put("/:id", userHandler.UpdateUser)        // Update user (self or admin only)
	users.Delete("/:id", userHandler.DeleteUser)     // Delete user (self or admin only)

	// Admin only routes
	users.Get("/", middleware.RequireRole("admin"), userHandler.GetAllUsers)  // List all users
} 