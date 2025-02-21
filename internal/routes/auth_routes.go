package routes

import (
	"example/API_Gateway/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, authHandler *handlers.AuthHandler) {
	// No middleware - public endpoint
	app.Get("/auth/tokens", authHandler.GetTokens)
}