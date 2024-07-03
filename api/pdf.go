package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/QuodEstDubitandum/iUtils/utils"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// /merge
// /split
// /encrypt
func HandlePDF(fileContext *utils.FileContext) error {
	if err := fileContext.Ctx.BodyParser(fileContext.Request); err != nil {
		log.Println("Could not parse request body: ", err)
		return errors.New("Invalid request body")
	}

	// if the inputformat sent by request dont match allowed ones
	if _, ok := fileContext.AllowedInputFormats[fileContext.Request.InputFormat]; !ok {
		fmt.Println("The allowed input seem to have been tampered with")
		return errors.New("The allowed input format seem to have been tampered with.")
	}

	// generate a unique hash as prefix for files so we have no collision
	fileContext.FilenamePrefix = utils.GenerateHash()

	switch fileContext.Ctx.Path() {
	case "/main_backend/merge":
		if err := mergePDF(fileContext); err != nil {
			return err
		}
	case "/main_backend/split":
		if err := splitPDF(fileContext); err != nil {
			return err
		}
	case "/main_backend/encrypt":
		fileContext.UserPassword = fileContext.Request.UserPW
		fileContext.AdminPassword = fileContext.Request.AdminPW
		if err := encryptPDF(fileContext); err != nil {
			return err
		}
	}

	err := utils.SendFileResponse(fileContext)
	return err
}

func mergePDF(fileContext *utils.FileContext) error {
	var totalFileSize int64
	var inputFilenameArray []string
	for _, file := range fileContext.Request.Files {
		totalFileSize += file.Size

		// if the actual file format is not the same as the information about it sent in request
		if file.ContentType != fileContext.Request.InputFormat {
			log.Println("File has incorrect input format")
			return errors.New(fmt.Sprintf("%s does not match the selected input format.", file.Name))
		}

		inputFilenameArray = append(inputFilenameArray, fmt.Sprintf("assets/%s/%s_%s", fileContext.Subfolder, fileContext.FilenamePrefix, file.Name))
	}

	// if combined file size is bigger than MaxFileSize
	if totalFileSize > fileContext.MaxFileSize {
		fmt.Println("Files too large")
		return errors.New(fmt.Sprintf("Your files are too large. Please select files with a combined maximum size of %dMB.", fileContext.MaxFileSize/1024/1024))
	}

	outputFilename := fmt.Sprintf("assets/%s/%s_merged_pdf", fileContext.Subfolder, fileContext.FilenamePrefix)

	for i, file := range fileContext.Request.Files {
		temp := i

		imageBytes, err := base64.StdEncoding.DecodeString(file.Data)
		if err != nil {
			log.Println("Could not decode base64 data into bytes: ", err)
			return errors.New("Could not decode base64 data into bytes")
		}
		err = os.WriteFile(inputFilenameArray[i], imageBytes, 0755)
		if err != nil {
			fmt.Println("Got error saving file to disk:", err)
			return errors.New("Couldnt save " + file.Name + " to disk.")
		}
		defer os.Remove(inputFilenameArray[temp])
	}

	err := api.MergeCreateFile(inputFilenameArray, outputFilename, nil)
	if err != nil {
		fmt.Println("PDF Merging failed: ", err)
		return errors.New("PDF Merging failed.")
	}
	fileContext.FilePaths = append(fileContext.FilePaths, outputFilename)
	return nil
}

func splitPDF(fileContext *utils.FileContext) error {
	file := fileContext.Request.Files[0]
	// if the actual file format is not the same as the information about it sent in request
	if file.ContentType != fileContext.Request.InputFormat {
		fmt.Println("File has incorrect input format")
		return errors.New(fmt.Sprintf("%s does not match the selected input format.", file.Name))
	}

	// if file size is bigger than MaxFileSize
	if file.Size > fileContext.MaxFileSize {
		fmt.Println("File too large")
		return errors.New(fmt.Sprintf("%s is too large. Please select files with a combined maximum size of %dMB.", file.Name, fileContext.MaxFileSize/1024/1024))
	}

	inputFilename := fmt.Sprintf("assets/%s/%s_%s", fileContext.Subfolder, fileContext.FilenamePrefix, file.Name)
	outputDir := fmt.Sprintf("assets/%s/%s", fileContext.Subfolder, fileContext.FilenamePrefix)

	err := os.Mkdir(outputDir, 0755)
	if err != nil {
		fmt.Println("Got error creating subfolder:", err)
		return errors.New("Couldnt create  " + outputDir)
	}

	imageBytes, err := base64.StdEncoding.DecodeString(file.Data)
	if err != nil {
		log.Println("Could not decode base64 data into bytes: ", err)
		return errors.New("Could not decode base64 data into bytes")
	}

	err = os.WriteFile(inputFilename, imageBytes, 0755)
	if err != nil {
		fmt.Println("Got error saving file to disk:", err)
		return errors.New("Couldnt save " + file.Name + " to disk.")
	}
	defer os.Remove(inputFilename)

	err = api.SplitFile(inputFilename, outputDir, 1, nil)
	if err != nil {
		fmt.Println("PDF Splitting failed: ", err)
		return errors.New("PDF Splitting failed.")
	}
	return nil
}

func encryptPDF(fileContext *utils.FileContext) error {
	var totalFileSize int64
	var inputFilenameArray []string
	var outputFilenameArray []string
	conf := model.NewAESConfiguration(fileContext.UserPassword, fileContext.AdminPassword, 256)

	for _, file := range fileContext.Request.Files {
		totalFileSize += file.Size

		// if the actual file format is not the same as the information about it sent in request
		if file.ContentType != fileContext.Request.InputFormat {
			fmt.Println("File has incorrect input format")
			return errors.New(fmt.Sprintf("%s does not match the selected input format.", file.Name))
		}

		inputFilenameArray = append(inputFilenameArray, fmt.Sprintf("assets/%s/%s_%s", fileContext.Subfolder, fileContext.FilenamePrefix, file.Name))
		outputFilenameArray = append(outputFilenameArray, fmt.Sprintf("assets/%s/%s_encrypted_%s", fileContext.Subfolder, fileContext.FilenamePrefix, file.Name))
	}

	// if combined file size is bigger than MaxFileSize
	if totalFileSize > fileContext.MaxFileSize {
		fmt.Println("Files too large")
		return errors.New(fmt.Sprintf("Your files are too large. Please select files with a combined maximum size of %dMB.", fileContext.MaxFileSize/1024/1024))
	}

	for i, file := range fileContext.Request.Files {
		temp := i

		imageBytes, err := base64.StdEncoding.DecodeString(file.Data)
		if err != nil {
			log.Println("Could not decode base64 data into bytes: ", err)
			return errors.New("Could not decode base64 data into bytes")
		}

		err = os.WriteFile(inputFilenameArray[temp], imageBytes, 0755)
		if err != nil {
			fmt.Println("Got error saving file to disk:", err)
			return errors.New("Couldnt save " + file.Name + " to disk.")
		}
		defer os.Remove(inputFilenameArray[temp])

		err = api.EncryptFile(inputFilenameArray[temp], outputFilenameArray[temp], conf)
		if err != nil {
			fmt.Println("PDF Encryption failed: ", err)
			return errors.New("PDF Encryption failed.")
		}
		fileContext.FilePaths = append(fileContext.FilePaths, outputFilenameArray[temp])
	}
	return nil
}
