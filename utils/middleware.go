package utils

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	var apiKey = os.Getenv("API_KEY")
	return func(c *fiber.Ctx) error {
		requestKey := c.Get("x-api-key")
		if requestKey != apiKey {
			log.Printf("Unauthorized access from: %s", c.IP())
			return c.Status(403).Send([]byte("Unauthorized"))
		}
		return c.Next()
	}
}
