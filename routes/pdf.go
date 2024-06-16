package routes

import (
	"github.com/QuodEstDubitandum/iUtils/api"
	"github.com/QuodEstDubitandum/iUtils/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterPDFRoutes(app *fiber.App) {
	pdfBackend := app.Group("/pdf")

	pdfBackend.Post("/merge", func(c *fiber.Ctx) error {
		allowedInputFormats := map[string]string{
			"application/pdf": ".pdf",
		}
		fileContext := &utils.FileContext{
			Ctx:                 c,
			AllowedInputFormats: allowedInputFormats,
			Subfolder:           "documents",
			MaxFileSize:         1024 * 1024 * 50,
		}
		return api.HandlePDF(fileContext)
	})

	pdfBackend.Post("/split", func(c *fiber.Ctx) error {
		allowedInputFormats := map[string]string{
			"application/pdf": ".pdf",
		}
		fileContext := &utils.FileContext{
			Ctx:                 c,
			AllowedInputFormats: allowedInputFormats,
			Subfolder:           "documents",
			MaxFileSize:         1024 * 1024 * 50,
		}
		return api.HandlePDF(fileContext)
	})

	pdfBackend.Post("/encrypt", func(c *fiber.Ctx) error {
		allowedInputFormats := map[string]string{
			"application/pdf": ".pdf",
		}
		fileContext := &utils.FileContext{
			Ctx:                 c,
			AllowedInputFormats: allowedInputFormats,
			Subfolder:           "documents",
			MaxFileSize:         1024 * 1024 * 50,
		}
		return api.HandlePDF(fileContext)
	})
}
