package utils

import (
	"github.com/gofiber/fiber/v2"
)

type FileContext struct {
	Ctx                  *fiber.Ctx
	Request              *RequestBody
	AllowedInputFormats  map[string]string
	AllowedOutputFormats map[string]string
	FilenamePrefix       string
	Subfolder            string
	MaxFileSize          int64
	UserPassword         string
	AdminPassword        string
	FilePaths            []string
}

type File struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"mimeType"`
	Data        string `json:"base64"`
}

type RequestBody struct {
	Files        []*File `json:"files"`
	InputFormat  string  `json:"input_format"`
	OutputFormat string  `json:"output_format"`
	UserPW       string  `json:"user_pw"`
	AdminPW      string  `json:"admin_pw"`
	Input        string  `json:"input"`
}
