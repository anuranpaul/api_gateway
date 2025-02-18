package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"example/API_Gateway/middleware"
)

func main() {
	// Load configuration
	config := LoadConfig()

	// Create new Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Global middleware
	app.Use(middleware.AuthMiddleware)

	// Setup routes
	setupRoutes(app, config)

	// Start server
	log.Printf("API Gateway running on port %s", config.Port)
	log.Fatal(app.Listen(":" + config.Port))
}

func setupRoutes(app *fiber.App, config *Config) {
	// Admin routes
	adminGroup := app.Group("/admin")
	adminGroup.Use(middleware.RequireRole("admin"))
	adminGroup.All("/*", ReverseProxy(config.AdminServiceURL))

	// User routes
	userGroup := app.Group("/users")
	userGroup.Use(middleware.RequireRole("user"))
	userGroup.All("/*", ReverseProxy(config.UserServiceURL))
}
