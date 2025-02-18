package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	// Create new Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Setup routes
	setupRoutes(app)

	// Start server
	log.Println("User Service running on port 5001")
	log.Fatal(app.Listen(":5001"))
}

func setupRoutes(app *fiber.App) {
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
