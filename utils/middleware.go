package utils

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

var apiKey = os.Getenv("API_KEY")

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestKey := c.GetReqHeaders()["x-api-key"]
		if len(requestKey) == 0 || requestKey[0] != apiKey {
			log.Printf("Expected apiKey %s, got %v", apiKey, requestKey)
			return c.Status(403).Send([]byte("Unauthorized"))
		}
		return c.Next()
	}
}
