package utils

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
)

type FileContext struct {
	Ctx                  *fiber.Ctx
	Files                []*multipart.FileHeader
	AllowedInputFormats  map[string]string
	AllowedOutputFormats map[string]string
	InputFormat          string
	OutputFormat         string
	FilenamePrefix       string
	Subfolder            string
	MaxFileSize          int64
	UserPassword         string
	AdminPassword        string
	FilePaths            []string
}

