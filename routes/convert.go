package routes

import (
	"github.com/QuodEstDubitandum/iUtils/api"
	"github.com/QuodEstDubitandum/iUtils/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterConvertRoutes(app *fiber.App) {
	convertBackend := app.Group("/convert")

	convertBackend.Post("/img-convert", func(c *fiber.Ctx) error {
		allowedInputFormats := map[string]string{
			"image/jpeg":    ".jpg",
			"image/png":     ".png",
			"image/svg+xml": ".svg",
			"image/webp":    ".webp",
			"image/heif":    ".heic",
		}
		allowedOutputFormats := map[string]string{
			"image/jpeg":      ".jpg",
			"image/png":       ".png",
			"image/svg+xml":   ".svg",
			"image/webp":      ".webp",
			"image/heif":      ".heic",
			"application/pdf": ".pdf",
		}
		fileContext := &utils.FileContext{
			Ctx:                  c,
			AllowedInputFormats:  allowedInputFormats,
			AllowedOutputFormats: allowedOutputFormats,
			Subfolder:            "images",
			MaxFileSize:          1024 * 1024 * 50,
		}
		return api.HandleImgConvert(fileContext)
	})

	convertBackend.Post("/mp4-gif", func(c *fiber.Ctx) error {
		fileContext := &utils.FileContext{
			Ctx: c,
			AllowedInputFormats: map[string]string{
				"video/mp4": ".mp4",
			},
			AllowedOutputFormats: map[string]string{
				"image/gif": ".gif",
			},
			Subfolder:   "videos",
			MaxFileSize: 1024 * 1024 * 50,
		}
		return api.HandleMP4ToGIF(fileContext)
	})

	convertBackend.Post("/mp4-mp3", func(c *fiber.Ctx) error {
		fileContext := &utils.FileContext{
			Ctx: c,
			AllowedInputFormats: map[string]string{
				"video/mp4": ".mp4",
			},
			AllowedOutputFormats: map[string]string{
				"audio/mpeg": ".mp3",
			},
			Subfolder:   "videos",
			MaxFileSize: 1024 * 1024 * 50,
		}
		return api.HandleMP4ToMP3(fileContext)
	})

	convertBackend.Post("/pdf-docx", func(c *fiber.Ctx) error {
		fileContext := &utils.FileContext{
			Ctx: c,
			AllowedInputFormats: map[string]string{
				"application/pdf": ".pdf",
			},
			AllowedOutputFormats: map[string]string{
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
			},
			Subfolder:   "documents",
			MaxFileSize: 1024 * 1024 * 50,
		}
		return api.HandlePdfToDocx(fileContext)
	})

	convertBackend.Post("/docx-pdf", func(c *fiber.Ctx) error {
		fileContext := &utils.FileContext{
			Ctx: c,
			AllowedInputFormats: map[string]string{
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
			},
			AllowedOutputFormats: map[string]string{
				"application/pdf": ".pdf",
			},
			Subfolder:   "documents",
			MaxFileSize: 1024 * 1024 * 50,
		}
		return api.HandleDocxToPdf(fileContext)
	})
}
