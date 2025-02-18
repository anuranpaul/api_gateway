package proxy

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// ReverseProxy forwards requests to the backend service
func ReverseProxy(target string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Forward the request to the target service
		if err := proxy.Do(c, target+c.Path()); err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error": "Proxy forwarding failed",
			})
		}
		return nil
	}
}
