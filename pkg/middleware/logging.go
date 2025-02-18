package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

// RequestLogger logs incoming requests using logrus
func RequestLogger(log *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Process request
		err := c.Next()
		
		// Calculate duration
		duration := time.Since(start)

		// Get status code
		statusCode := c.Response().StatusCode()

		// Create log entry
		log.WithFields(logrus.Fields{
			"method":     c.Method(),
			"path":       c.Path(),
			"ip":         c.IP(),
			"status":     statusCode,
			"duration":   duration,
			"user_agent": c.Get("User-Agent"),
			"token":      c.Get("Authorization"),
		}).Info("Request processed")

		return err
	}
}

// InitLogger initializes and configures logrus
func InitLogger() *logrus.Logger {
	log := logrus.New()
	
	// Configure logrus
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	
	// Set log level
	log.SetLevel(logrus.InfoLevel)
	
	return log
}
