package api

import (
	"errors"
	"fmt"
	"os"

	"github.com/QuodEstDubitandum/iUtils/utils"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// /merge
// /split
// /encrypt
func HandlePDF(fileContext *utils.FileContext) error {
	form, err := fileContext.Ctx.MultipartForm()
	if err == nil {
		fileContext.Files = form.File["files"]
		fileContext.InputFormat = form.Value["input_format"][0]

		// if the inputformat sent by request dont match allowed ones
		if _, ok := fileContext.AllowedInputFormats[fileContext.InputFormat]; !ok {
			fmt.Println("The allowed input seem to have been tampered with")
			return errors.New("The allowed input format seem to have been tampered with.")
		}

		// generate a unique hash as prefix for files so we have no collision
		fileContext.FilenamePrefix = utils.GenerateHash()

		switch fileContext.Ctx.Path() {
		case "/main_backend/merge":
			err = mergePDF(fileContext)
		case "/main_backend/split":
			err = splitPDF(fileContext)
		case "/main_backend/encrypt":
			fileContext.UserPassword = form.Value["user_pw"][0]
			fileContext.AdminPassword = form.Value["admin_pw"][0]
			err = encryptPDF(fileContext)
		}
		if err != nil {
			return err
		}

		err := utils.SendFileResponse(fileContext)
		return err
	}
	return errors.New("The allowed content type seems to have been tampered with.")
}

func mergePDF(fileContext *utils.FileContext) error {
	var totalFileSize int64
	var inputFilenameArray []string
	for _, file := range fileContext.Files {
		totalFileSize += file.Size

		// if the actual file format is not the same as the information about it sent in request
		if file.Header["Content-Type"][0] != fileContext.InputFormat {
			fmt.Println("File has incorrect input format")
			return errors.New(fmt.Sprintf("%s does not match the selected input format.", file.Filename))
		}

		inputFilenameArray = append(inputFilenameArray, fmt.Sprintf("assets/%s/%s_%s", fileContext.Subfolder, fileContext.FilenamePrefix, file.Filename))
	}

	// if combined file size is bigger than MaxFileSize
	if totalFileSize > fileContext.MaxFileSize {
		fmt.Println("Files too large")
		return errors.New(fmt.Sprintf("Your files are too large. Please select files with a combined maximum size of %dMB.", fileContext.MaxFileSize/1024/1024))
	}

	outputFilename := fmt.Sprintf("assets/%s/%s_merged_pdf", fileContext.Subfolder, fileContext.FilenamePrefix)

	for i, file := range fileContext.Files {
		temp := i
		err := fileContext.Ctx.SaveFile(file, inputFilenameArray[i])
		if err != nil {
			fmt.Println("Got error saving file to disk:", err)
			return errors.New("Couldnt save " + file.Filename + " to disk.")
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
	file := fileContext.Files[0]
	// if the actual file format is not the same as the information about it sent in request
	if file.Header["Content-Type"][0] != fileContext.InputFormat {
		fmt.Println("File has incorrect input format")
		return errors.New(fmt.Sprintf("%s does not match the selected input format.", file.Filename))
	}

	// if file size is bigger than MaxFileSize
	if file.Size > fileContext.MaxFileSize {
		fmt.Println("File too large")
		return errors.New(fmt.Sprintf("%s is too large. Please select files with a combined maximum size of %dMB.", file.Filename, fileContext.MaxFileSize/1024/1024))
	}

	inputFilename := fmt.Sprintf("assets/%s/%s_%s", fileContext.Subfolder, fileContext.FilenamePrefix, file.Filename)
	outputDir := fmt.Sprintf("assets/%s/%s", fileContext.Subfolder, fileContext.FilenamePrefix)

	err := os.Mkdir(outputDir, 0755)
	if err != nil {
		fmt.Println("Got error creating subfolder:", err)
		return errors.New("Couldnt create  " + outputDir)
	}
	err = fileContext.Ctx.SaveFile(file, inputFilename)
	if err != nil {
		fmt.Println("Got error saving file to disk:", err)
		return errors.New("Couldnt save " + file.Filename + " to disk.")
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

	for _, file := range fileContext.Files {
		totalFileSize += file.Size

		// if the actual file format is not the same as the information about it sent in request
		if file.Header["Content-Type"][0] != fileContext.InputFormat {
			fmt.Println("File has incorrect input format")
			return errors.New(fmt.Sprintf("%s does not match the selected input format.", file.Filename))
		}

		inputFilenameArray = append(inputFilenameArray, fmt.Sprintf("assets/%s/%s_%s", fileContext.Subfolder, fileContext.FilenamePrefix, file.Filename))
		outputFilenameArray = append(outputFilenameArray, fmt.Sprintf("assets/%s/%s_encrypted_%s", fileContext.Subfolder, fileContext.FilenamePrefix, file.Filename))
	}

	// if combined file size is bigger than MaxFileSize
	if totalFileSize > fileContext.MaxFileSize {
		fmt.Println("Files too large")
		return errors.New(fmt.Sprintf("Your files are too large. Please select files with a combined maximum size of %dMB.", fileContext.MaxFileSize/1024/1024))
	}

	for i, file := range fileContext.Files {
		temp := i
		err := fileContext.Ctx.SaveFile(file, inputFilenameArray[temp])
		if err != nil {
			fmt.Println("Got error saving file to disk:", err)
			return errors.New("Couldnt save " + file.Filename + " to disk.")
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
