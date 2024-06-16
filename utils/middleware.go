package utils

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var apiKey = os.Getenv("API-KEY")

func AuthMiddleware(redis *redis.Client, ctx context.Context) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestKey := c.GetReqHeaders()["x-api-key"][0]
		if requestKey != apiKey {
			return c.Status(403).Send([]byte("Unauthorized"))
		}
		return c.Next()
	}
}

