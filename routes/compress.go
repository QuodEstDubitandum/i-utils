package routes

import (
	"github.com/QuodEstDubitandum/iUtils/api"
	"github.com/QuodEstDubitandum/iUtils/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterCompressRoutes(app *fiber.App) {
	compressBackend := app.Group("/compress")

	compressBackend.Post("/jpg", func(c *fiber.Ctx) error {
		allowedInputFormats := map[string]string{
			"image/jpeg": ".jpg",
			"image/webp": ".webp",
		}
		fileContext := &utils.FileContext{
			Ctx:                 c,
			AllowedInputFormats: allowedInputFormats,
			Subfolder:           "images",
			MaxFileSize:         1024 * 1024 * 10,
		}
		return api.HandleImgCompression(fileContext)
	})

	compressBackend.Post("/png", func(c *fiber.Ctx) error {
		allowedInputFormats := map[string]string{
			"image/png": ".png",
		}
		fileContext := &utils.FileContext{
			Ctx:                 c,
			AllowedInputFormats: allowedInputFormats,
			Subfolder:           "images",
			MaxFileSize:         1024 * 1024 * 10,
		}
		return api.HandleImgCompression(fileContext)
	})

	compressBackend.Post("/pdf", func(c *fiber.Ctx) error {
		allowedInputFormats := map[string]string{
			"application/pdf": ".pdf",
		}
		fileContext := &utils.FileContext{
			Ctx:                 c,
			AllowedInputFormats: allowedInputFormats,
			Subfolder:           "documents",
			MaxFileSize:         1024 * 1024 * 50,
		}
		return api.HandleImgCompression(fileContext)
	})
}
