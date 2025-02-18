package services

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// Admin represents the structure for the admin
type Admin struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func StartAdminService() {
	// Create new Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Setup routes
	setupAdminRoutes(app)

	// Start server
	log.Println("Admin Service running on port 5003")
	log.Fatal(app.Listen(":5003"))
}

func setupAdminRoutes(app *fiber.App) {
	// Admin handlers
	app.Get("/admin", getAdminHandler)
	app.Get("/admin/dashboard", getAdminHandler)
}

func getAdminHandler(c *fiber.Ctx) error {
	admin := Admin{
		ID:   1,
		Name: "Admin User",
	}

	return c.JSON(admin)
}
