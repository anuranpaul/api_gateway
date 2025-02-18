package main

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/sirupsen/logrus"
	"example/API_Gateway/middleware"
	"example/API_Gateway/metrics"
)

var logger *logrus.Logger

func main() {
	// Initialize logger
	logger = middleware.InitLogger()

	// Load configuration
	config := LoadConfig()
	logger.Info("Configuration loaded successfully")

	// Create new Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.WithError(err).Error("Application error occurred")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Apply middleware
	app.Use(middleware.RequestLogger(logger))
	setupTokenRateLimiter(app)
	setupIPRateLimiter(app)
	app.Use(middleware.AuthMiddleware)
	app.Use(metrics.PrometheusMiddleware)

	// Setup routes
	setupRoutes(app, config)
	logger.Info("Routes configured successfully")

	// Add metrics endpoint
	app.Get("/metrics", metrics.MetricsHandler())

	// Start server
	logger.WithFields(logrus.Fields{
		"port": config.Port,
	}).Info("API Gateway starting")
	
	if err := app.Listen(":" + config.Port); err != nil {
		logger.WithError(err).Fatal("Server failed to start")
	}
}

// setupTokenRateLimiter configures and applies token-based rate limiting
func setupTokenRateLimiter(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        20,
		Expiration: 30 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			token := c.Get("Authorization")
			logger.WithFields(logrus.Fields{
				"token": token,
			}).Debug("Token rate limit check")
			return "token:" + token
		},
		LimitReached: func(c *fiber.Ctx) error {
			logger.WithFields(logrus.Fields{
				"token": c.Get("Authorization"),
			}).Warn("Token rate limit exceeded")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Token rate limit exceeded. Please try again after 30 seconds.",
			})
		},
		Storage: limiter.ConfigDefault.Storage,
		LimiterMiddleware: limiter.FixedWindow{},
		Next: func(c *fiber.Ctx) bool {
			return c.Get("Authorization") == ""
		},
	}))
}

// setupIPRateLimiter configures and applies IP-based rate limiting
func setupIPRateLimiter(app *fiber.App) {
	app.Use(limiter.New(limiter.Config{
		Max:        30,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.IP()
			logger.WithFields(logrus.Fields{
				"ip": ip,
			}).Debug("IP rate limit check")
			return "ip:" + ip
		},
		LimitReached: func(c *fiber.Ctx) error {
			logger.WithFields(logrus.Fields{
				"ip": c.IP(),
			}).Warn("IP rate limit exceeded")
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
	// Admin routes
	adminGroup := app.Group("/admin")
	adminGroup.Use(middleware.RequireRole("admin"))
	adminGroup.All("/*", ReverseProxy(config.AdminServiceURL))
	logger.Info("Admin routes configured")

	// User routes
	userGroup := app.Group("/users")
	userGroup.Use(middleware.RequireRole("user"))
	userGroup.All("/*", ReverseProxy(config.UserServiceURL))
	logger.Info("User routes configured")
}
