package utils

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

var apiKey = os.Getenv("API-KEY")

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestKey := c.GetReqHeaders()["x-api-key"][0]
		if requestKey != apiKey {
			return c.Status(403).Send([]byte("Unauthorized"))
		}
		return c.Next()
	}
}
