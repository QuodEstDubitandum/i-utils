package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// send the files to the frontend and remove the output file from disk afterwards
func SendFileResponse(fileContext *FileContext) error {
	// remove the output files after sending response
	for _, filepath := range fileContext.FilePaths {
		defer os.Remove(filepath)
	}

	// if we split pdfs we want to add the filepaths like this
	if fileContext.Ctx.Path() == "/pdf/split" {
		rootDir := fmt.Sprintf("assets/%s/%s", fileContext.Subfolder, fileContext.FilenamePrefix)
		err := filepath.Walk(rootDir, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing path %q: %v\n", filePath, err)
				return err
			}
			if filePath != rootDir {
				fileContext.FilePaths = append(fileContext.FilePaths, filePath)
			}
			return nil
		})
		if err != nil {
			return err
		}
		defer os.RemoveAll(rootDir)
	}

	// if only a single file
	if len(fileContext.FilePaths) == 1 {
		fileContext.Ctx.SendFile(fileContext.FilePaths[0])
		return nil
	}

	// if multiple files
	zipName, err := zipFiles(fileContext)
	if err != nil {
		log.Println(err)
		os.Remove(zipName)
		return err
	}
	fileContext.Ctx.Set("Content-Type", "application/zip")
	fileContext.Ctx.SendFile(zipName)
	os.Remove(zipName)
	return nil
}

func zipFiles(fileContext *FileContext) (string, error) {
	// create a new zip archive
	zipName := fmt.Sprintf("assets/%s/%s_archive.zip", fileContext.Subfolder, fileContext.FilenamePrefix)
	zipFile, err := os.Create(zipName)
	if err != nil {
		log.Println("Couldnt create zip archive: ", err)
		return zipName, errors.New("Something went wrong on the server.")
	}
	defer zipFile.Close()

	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	// iterate over all the files and add them to the zip archive
	for _, filepath := range fileContext.FilePaths {
		err := addFileToZip(writer, filepath, fileContext)
		if err != nil {
			log.Println(err)
			return zipName, err
		}
	}
	return zipName, nil
}

func addFileToZip(writer *zip.Writer, filepath string, fileContext *FileContext) error {
	fileToZip, err := os.Open(filepath)
	if err != nil {
		log.Println("Couldnt open the file for zipping: ", err)
		return errors.New("Something went wrong on the server.")
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		log.Println("Couldnt get the file information for zipping: ", err)
		return errors.New("Something went wrong on the server.")
	}

	// Create a zip header for the file
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		log.Println("Couldnt create zip header: ", err)
		return errors.New("Something went wrong on the server.")
	}
	if fileContext.Ctx.Path() == "/pdf/split" {
		header.Name = strings.ReplaceAll(filepath, fmt.Sprintf("assets/%s/%s/%s_", fileContext.Subfolder, fileContext.FilenamePrefix, fileContext.FilenamePrefix), "")
	} else {
		header.Name = strings.ReplaceAll(filepath, fmt.Sprintf("assets/%s/%s_", fileContext.Subfolder, fileContext.FilenamePrefix), "")
	}

	// Add the header to the zip writer
	w, err := writer.CreateHeader(header)
	if err != nil {
		log.Println("Couldnt add header to zip:  ", err)
		return errors.New("Something went wrong on the server.")
	}

	// Write the file content to the zip archive
	_, err = io.Copy(w, fileToZip)
	if err != nil {
		log.Println("Couldnt write file to zip: ", err)
		return errors.New("Something went wrong on the server.")
	}
	return nil
}
