package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func ConvertFile(file *File, fileContext *FileContext) error {
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
	fileNameArray := strings.Split(inputFilename, ".")
	outputFilename := strings.Join(fileNameArray[:len(fileNameArray)-1], ".") + fileContext.AllowedOutputFormats[fileContext.Request.OutputFormat]

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
	case "/convert/mp4-gif":
		return convertMP4ToGif(fileContext, inputFilename, outputFilename)
	case "/convert/mp4-mp3":
		return convertMP4ToMP3(fileContext, inputFilename, outputFilename)
	case "/convert/pdf-docx":
		return convertPDF(fileContext, inputFilename, outputFilename)
	case "/convert/docx-pdf":
		return convertDocx(fileContext, inputFilename, outputFilename)
	default:
		return convertImage(fileContext, inputFilename, outputFilename)
	}
}

func convertImage(fileContext *FileContext, inputFilename string, outputFilename string) error {
	cmd := exec.Command("convert", inputFilename, outputFilename)
	fmt.Println(inputFilename, outputFilename)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println("Image conversion failed: ", err)
		log.Println("Convert output: ", string(output))
		return errors.New("Image conversion failed.")
	}
	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func convertMP4ToMP3(fileContext *FileContext, inputFilename string, outputFilename string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFilename, "-q:a", "0", "-map", "a", outputFilename)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return errors.New("Video conversion failed.")
	}
	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func convertMP4ToGif(fileContext *FileContext, inputFilename string, outputFilename string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFilename, "-ss", "00:00:00.000", "-t", "00:00:07.000", "-vf", "fps=24,scale=600:-1:flags=lanczos", "-c:v", "gif", outputFilename)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return errors.New("Video conversion failed.")
	}
	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func convertPDF(fileContext *FileContext, inputFilename string, outputFilename string) error {
	cmd := exec.Command("libreoffice", "--infilter=writer_pdf_import", "--convert-to", "docx", "--outdir", fmt.Sprintf("assets/%s/", fileContext.Subfolder), inputFilename)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return errors.New("PDF document converting failed.")
	}

	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func convertDocx(fileContext *FileContext, inputFilename string, outputFilename string) error {
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf", "--outdir", fmt.Sprintf("assets/%s/", fileContext.Subfolder), inputFilename)
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return errors.New("Word document converting failed.")
	}

	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}
