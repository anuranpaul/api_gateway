package main

import (
	"example/API_Gateway/internal/config"
	"example/API_Gateway/pkg/metrics"
	"example/API_Gateway/pkg/middleware"
	"example/API_Gateway/pkg/proxy"
	"time"

	"example/API_Gateway/internal/db"
	"example/API_Gateway/internal/handlers"
	"example/API_Gateway/internal/repository"
	"example/API_Gateway/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/sirupsen/logrus"

	"example/API_Gateway/internal/cache"
)

var logger *logrus.Logger

func main() {
	// Initialize logger
	logger = middleware.InitLogger()

	// Load configuration
	config := config.LoadConfig()
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
	redisClient := cache.NewRedisClient()
	app.Use(middleware.AuthMiddleware(redisClient))
	app.Use(metrics.PrometheusMiddleware)

	// Initialize database
	database, err := db.NewDB(config.DatabaseURL)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer database.Close()

	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(database)
	userHandler := handlers.NewUserHandler(userRepo)
	authHandler := handlers.NewAuthHandler(redisClient)

	// Setup routes
	setupRoutes(app, config)
	routes.SetupUserRoutes(app, userHandler)
	routes.SetupAuthRoutes(app, authHandler)
	logger.Info("Routes configured successfully")

	// Add metrics endpoint with admin protection
	app.Get("/metrics", middleware.RequireRole("admin"), metrics.MetricsHandler())

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
func setupRoutes(app *fiber.App, config *config.Config) {
	// Admin routes
	adminGroup := app.Group("/admin")
	adminGroup.Use(middleware.RequireRole("admin"))
	adminGroup.All("/*", proxy.ReverseProxy(config.AdminServiceURL))
	logger.Info("Admin routes configured")

	// User routes
	userGroup := app.Group("/users")
	userGroup.Use(middleware.RequireRole("user"))
	userGroup.All("/*", proxy.ReverseProxy(config.UserServiceURL))
	logger.Info("User routes configured")
}