package routes

import (
	"encoding/base64"
	"errors"
	"log"
	"strconv"

	"github.com/QuodEstDubitandum/iUtils/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterHashRoutes(app *fiber.App) {
	hashBackend := app.Group("/hash")

	hashBackend.Post("/encode-base64", func(c *fiber.Ctx) error {
		requestBody := &utils.RequestBody{}
		if err := c.BodyParser(requestBody); err != nil {
			log.Println("Could not parse request body: ", err)
			return errors.New("Invalid request body")
		}

		encodedString := base64.StdEncoding.EncodeToString([]byte(requestBody.Input))
		c.Status(200).Send([]byte(encodedString))
		return nil
	})

	hashBackend.Post("/decode-base64", func(c *fiber.Ctx) error {
		requestBody := &utils.RequestBody{}
		if err := c.BodyParser(requestBody); err != nil {
			log.Println("Could not parse request body: ", err)
			return errors.New("Invalid request body")
		}

		decodedString, err := base64.StdEncoding.DecodeString(requestBody.Input)
		if err != nil {
			log.Println("Could not parse the string from base64: ", err)
			return errors.New("Invalid base64 data")
		}

		c.Status(200).Send(decodedString)
		return nil
	})

	hashBackend.Post("/decode-ascii", func(c *fiber.Ctx) error {
		requestBody := &utils.RequestBody{}
		if err := c.BodyParser(requestBody); err != nil {
			log.Println("Could not parse request body: ", err)
			return errors.New("Invalid request body")
		}

		ascii, err := strconv.Atoi(requestBody.Input)
		if err != nil {
			log.Println("Invalid ascii character provided: ", err)
			return errors.New("Invalid ascii code")
		}

		if ascii < 0 || ascii > 255 {
			log.Println("Invalid ascii character provided: ", ascii)
			return errors.New("Only codes between 0 and 255 are valid ascii codes")
		}

		c.Status(200).Send([]byte(string(rune(ascii))))
		return nil
	})

	hashBackend.Post("/encode-hex", func(c *fiber.Ctx) error { return nil })
	hashBackend.Post("/decode-hex", func(c *fiber.Ctx) error { return nil })

	hashBackend.Post("/encode-binary", func(c *fiber.Ctx) error { return nil })
	hashBackend.Post("/decode-binary", func(c *fiber.Ctx) error { return nil })
}
