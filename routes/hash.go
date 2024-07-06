package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

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

		if ascii < 33 || ascii > 126 {
			log.Println("Invalid ascii character provided: ", ascii)
			return errors.New("Only codes between 0 and 255 are valid ascii codes")
		}

		c.Status(200).Send([]byte(string(rune(ascii))))
		return nil
	})

	hashBackend.Post("/encode-unicode", func(c *fiber.Ctx) error {
		requestBody := &utils.RequestBody{}
		if err := c.BodyParser(requestBody); err != nil {
			log.Println("Could not parse request body: ", err)
			return errors.New("Invalid request body")
		}

		if len(requestBody.Input) < 1 || len(requestBody.Input) > 4 {
			log.Println("Incorrect unicode length: ", len(requestBody.Input))
			return errors.New(fmt.Sprintf("Invalid length of utf-8 char: %d", len(requestBody.Input)))
		}

		r := []rune(requestBody.Input)[0]

		c.Status(200).Send([]byte(fmt.Sprintf("U+%04X", r)))
		return nil
	})

	hashBackend.Post("/decode-unicode", func(c *fiber.Ctx) error {
		requestBody := &utils.RequestBody{}
		if err := c.BodyParser(requestBody); err != nil {
			log.Println("Could not parse request body: ", err)
			return errors.New("Invalid request body")
		}

		cutString, bool := strings.CutPrefix(requestBody.Input, "U+")
		if !bool {
			log.Println("Invalid unicode: ", requestBody.Input)
			return errors.New("Invalid unicode character")
		}

		unicodeHex, err := strconv.ParseInt(cutString, 16, 32)
		if err != nil {
			log.Println("Invalid hex code: ", cutString)
			return errors.New("Invalid unicode character")
		}

		c.Status(200).Send([]byte(string(rune(unicodeHex))))
		return nil
	})

	hashBackend.Post("/sha256", func(c *fiber.Ctx) error {
		requestBody := &utils.RequestBody{}
		if err := c.BodyParser(requestBody); err != nil {
			log.Println("Could not parse request body: ", err)
			return errors.New("Invalid request body")
		}

		sha256Hash := sha256.New()
		sha256Hash.Write([]byte(requestBody.Input))
		sha256Sum := sha256Hash.Sum(nil)

		c.Status(200).Send(sha256Sum)
		return nil
	})

	hashBackend.Post("/hmac-sha256", func(c *fiber.Ctx) error {
		requestBody := &utils.RequestBody{}
		if err := c.BodyParser(requestBody); err != nil {
			log.Println("Could not parse request body: ", err)
			return errors.New("Invalid request body")
		}

		hmacSha256 := hmac.New(sha256.New, []byte(requestBody.Secret))
		hmacSha256.Write([]byte(requestBody.Input))
		hmacSha256Sum := hmacSha256.Sum(nil)

		c.Status(200).Send(hmacSha256Sum)
		return nil
	})
}
