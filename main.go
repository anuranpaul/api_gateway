package main

import (
	"fmt"
	"log"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"example/API_Gateway/middleware"
	"example/API_Gateway/metrics"
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

	// Apply Global Rate Limiting
	setupTokenRateLimiter(app)  // Token-based rate limiting
	setupIPRateLimiter(app)     // IP-based rate limiting
	
	// Global middleware
	app.Use(middleware.AuthMiddleware)

	// Setup routes
	setupRoutes(app, config)

	// Start server
	log.Printf("API Gateway running on port %s", config.Port)
	log.Fatal(app.Listen(":" + config.Port))
}

// setupTokenRateLimiter configures and applies token-based rate limiting
func setupTokenRateLimiter(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			token := c.Get("Authorization")
			fmt.Println("Token rate limiting for:", token)
			return "token:" + token // Prefix to avoid conflicts with IP limiter
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Token rate limit exceeded. Please try again after 30 seconds.",
			})
		},
		Storage: limiter.ConfigDefault.Storage,
		LimiterMiddleware: limiter.FixedWindow{},
		Next: func(c *fiber.Ctx) bool {
			return c.Get("Authorization") == "" // Skip if no token (let IP limiter handle it)
		},
	}))
}

// setupIPRateLimiter configures and applies IP-based rate limiting
func setupIPRateLimiter(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        30,         // Higher limit for IP-based
		Expiration: 60 * time.Second, // Longer window for IP-based
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.IP()
			fmt.Println("IP rate limiting for:", ip)
			return "ip:" + ip 
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "IP rate limit exceeded. Please try again after 60 seconds.",
			})
		},
		Storage: limiter.ConfigDefault.Storage,
		LimiterMiddleware: limiter.FixedWindow{},
		Next: func(c *fiber.Ctx) bool {
			return false 
		},
	}))
}

// setupRoutes configures all the routes for the application
func setupRoutes(app *fiber.App, config *Config) {
	// Admin routes with role-based access
	adminGroup := app.Group("/admin")
	adminGroup.Use(middleware.RequireRole("admin"))
	adminGroup.All("/*", ReverseProxy(config.AdminServiceURL))

	// User routes with role-based access
	userGroup := app.Group("/users")
	userGroup.Use(middleware.RequireRole("user"))
	userGroup.All("/*", ReverseProxy(config.UserServiceURL))
}
