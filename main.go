package main

import (
	"fmt"
	"log"
	"os"

	"github.com/QuodEstDubitandum/iUtils/routes"
	"github.com/QuodEstDubitandum/iUtils/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 200 * 1024 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("REQUEST_URL"),
		AllowHeaders: "Origin, Content-Type, Accept, X-API-Key",
		AllowMethods: "POST, OPTIONS, GET",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).Send([]byte("OK"))
	})

	if os.Getenv("ENVIRONMENT") == "prod" {
		app.Use(utils.AuthMiddleware())
	}

	routes.RegisterConvertRoutes(app)
	routes.RegisterCompressRoutes(app)
	routes.RegisterPDFRoutes(app)
	routes.RegisterHashRoutes(app)

	port := os.Getenv("PORT")
	err = app.ListenTLS(fmt.Sprintf(":%s", port), "./certs/cert.pem", "./certs/key.pem")
	if err != nil {
		log.Printf("Couldnt start server on port %s: %v", port, err)
		return
	}
	log.Printf("Listening on port %s", port)
}
