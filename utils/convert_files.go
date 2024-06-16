package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"strings"
)

func ConvertFile(file *multipart.FileHeader, fileContext *FileContext) error{
	// if file size is bigger than MaxFileSize
	if file.Size > fileContext.MaxFileSize{
		fmt.Println("File too large")
		return errors.New(fmt.Sprintf("%s is too large. Please select files with a maximum size of %dMB.", file.Filename, fileContext.MaxFileSize/1024/1024))
	}
	// if the actual file format is not the same as the information about it sent in request
	if file.Header["Content-Type"][0] != fileContext.InputFormat{
		fmt.Println("File has incorrect input format")
		return errors.New(fmt.Sprintf("%s does not match the selected input format.", file.Filename))
	}

	inputFilename := fmt.Sprintf("assets/%s/%s_%s", fileContext.Subfolder ,fileContext.FilenamePrefix, file.Filename)
	fileNameArray := strings.Split(inputFilename, ".")
	outputFilename := strings.Join(fileNameArray[:len(fileNameArray)-1], ".") + fileContext.AllowedOutputFormats[fileContext.OutputFormat]

	err := fileContext.Ctx.SaveFile(file, inputFilename)
	if err != nil{
		fmt.Println("Got error saving file to disk:", err)
		return errors.New("Couldnt save " + file.Filename + " to disk.")
	}
	defer os.Remove(inputFilename)

	// convert input file and save output to disk
	switch fileContext.Ctx.Path(){
		case "/convert_backend/mp4-gif":
			return convertMP4ToGif(fileContext, inputFilename, outputFilename)
		case "/convert_backend/mp4-mp3":
			return convertMP4ToMP3(fileContext, inputFilename, outputFilename)
		case "/convert_backend/pdf-docx":
			return convertPDF(fileContext, inputFilename, outputFilename)
		case "/convert_backend/docx-pdf":
			return convertDocx(fileContext, inputFilename, outputFilename)
		default:
			return convertImage(fileContext, inputFilename, outputFilename)
	}
}

func convertImage(fileContext *FileContext, inputFilename string, outputFilename string) error{
	cmd := exec.Command("convert", inputFilename, outputFilename)
	if err := cmd.Run(); err != nil{
		fmt.Println("Image conversion failed: ", err)
		return errors.New("Image conversion failed.")
	}
	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func convertMP4ToMP3(fileContext *FileContext, inputFilename string, outputFilename string) error{
	cmd := exec.Command("ffmpeg", "-i", inputFilename, "-q:a", "0", "-map", "a", outputFilename)
	if err := cmd.Run(); err != nil{
		fmt.Println(err)
		return errors.New("Video conversion failed.")
	}
	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func convertMP4ToGif(fileContext *FileContext, inputFilename string, outputFilename string) error{
	cmd := exec.Command("ffmpeg", "-i", inputFilename, "-ss", "00:00:00.000", "-t", "00:00:07.000", "-vf", "fps=24,scale=600:-1:flags=lanczos", "-c:v", "gif", outputFilename)
	if err := cmd.Run(); err != nil{
		fmt.Println(err)
		return errors.New("Video conversion failed.")
	}
	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func convertPDF(fileContext *FileContext, inputFilename string, outputFilename string) error{
	cmd := exec.Command("libreoffice", "--infilter=writer_pdf_import", "--convert-to", "docx","--outdir", fmt.Sprintf("assets/%s/", fileContext.Subfolder), inputFilename)
	if err := cmd.Run(); err != nil{
		fmt.Println(err)
		return errors.New("PDF document converting failed.")
	}

	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func convertDocx(fileContext *FileContext, inputFilename string, outputFilename string) error{
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf","--outdir", fmt.Sprintf("assets/%s/", fileContext.Subfolder), inputFilename)
	if err := cmd.Run(); err != nil{
		fmt.Println(err)
		return errors.New("Word document converting failed.")
	}

	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}