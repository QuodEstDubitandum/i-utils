package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func CompressFile(file *File, fileContext *FileContext) error {
	// if file size is bigger than MaxFileSize
	if file.Size > fileContext.MaxFileSize {
		log.Println("File too large")
		return errors.New(fmt.Sprintf("%s is too large. Please select files with a maximum size of %dMB.", file.Name, fileContext.MaxFileSize/1024/1024))
	}
	// if the actual file format is not the same as the information about it sent in request
	if file.ContentType != fileContext.Request.InputFormat {
		log.Println("File has incorrect input format")
		return errors.New(fmt.Sprintf("%s does not match the selected input format.", file.Name))
	}

	inputFilename := fmt.Sprintf("assets/%s/%s_%s", fileContext.Subfolder, fileContext.FilenamePrefix, file.Name)
	outputFilename := fmt.Sprintf("assets/%s/%s_compressed-%s", fileContext.Subfolder, fileContext.FilenamePrefix, file.Name)

	imageBytes, err := base64.StdEncoding.DecodeString(file.Data)
	if err != nil {
		log.Println("Could not decode base64 data into bytes: ", err)
		return errors.New("Could not decode base64 data into bytes")
	}

	err = os.WriteFile(inputFilename, imageBytes, 0755)
	if err != nil {
		log.Println("Got error saving file to disk:", err)
		return errors.New("Couldnt save " + file.Name + " to disk.")
	}
	defer os.Remove(inputFilename)

	// convert input file and save output to disk
	switch fileContext.Ctx.Path() {
	case "/compress/pdf":
		return compressPdf(fileContext, inputFilename, outputFilename)
	case "/compress/png":
		return compressPng(fileContext, inputFilename, outputFilename)
	default:
		return compressJpg(fileContext, inputFilename, outputFilename)
	}
}

func compressJpg(fileContext *FileContext, inputFilename string, outputFilename string) error {
	cmd := exec.Command("convert", inputFilename, "-quality", "50", outputFilename)
	if err := cmd.Run(); err != nil {
		log.Println("Image compression failed: ", err)
		return errors.New("Image compression failed.")
	}
	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func compressPng(fileContext *FileContext, inputFilename string, outputFilename string) error {
	cmd := exec.Command("pngquant", inputFilename, "-o", outputFilename)
	if err := cmd.Run(); err != nil {
		log.Println("Image compression failed: ", err)
		return errors.New("Image compression failed.")
	}
	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func compressPdf(fileContext *FileContext, inputFilename string, outputFilename string) error {
	err := api.OptimizeFile(inputFilename, outputFilename, nil)

	if err != nil {
		log.Println("PDF compression failed: ", err)
		return errors.New("PDF compression failed.")
	}
	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}
