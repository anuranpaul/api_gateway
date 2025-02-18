package services

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func StartUserService() {
	// Create new Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Setup routes
	setupUserRoutes(app)

	// Start server
	log.Println("User Service running on port 5001")
	log.Fatal(app.Listen(":5001"))
}

func setupUserRoutes(app *fiber.App) {
	// User handlers
	app.Get("/users", getUserHandler)
	app.Get("/users/test", getUserHandler)
}

func getUserHandler(c *fiber.Ctx) error {
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}

	return c.JSON(user)
}
