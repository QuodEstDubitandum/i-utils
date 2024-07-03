package utils

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

var apiKey = os.Getenv("API_KEY")

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestKey := c.GetReqHeaders()["x-api-key"]
		if len(requestKey) == 0 || requestKey[0] != apiKey {
			return c.Status(403).Send([]byte("Unauthorized"))
		}
		return c.Next()
	}
}
