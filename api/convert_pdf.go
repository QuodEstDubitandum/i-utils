package api

import (
	"errors"
	"log"
	"sync"

	"github.com/QuodEstDubitandum/iUtils/utils"
)

// /pdf-docx
func HandlePdfToDocx(fileContext *utils.FileContext) error {
	if err := fileContext.Ctx.BodyParser(fileContext.Request); err != nil {
		log.Println("Could not parse request body: ", err)
		return errors.New("Invalid request body")
	}
	// generate a unique hash as prefix for files so we have no collision
	fileContext.FilenamePrefix = utils.GenerateHash()

	// convert the files in parallel using go routines
	errorChan := make(chan error, len(fileContext.Request.Files))
	var wg sync.WaitGroup

	for _, file := range fileContext.Request.Files {
		wg.Add(1)
		go func(file *utils.File) {
			defer wg.Done()
			err := utils.ConvertFile(file, fileContext)
			if err != nil {
				errorChan <- err
			}
		}(file)
	}
	wg.Wait()

	// return err if any of the file convertings threw an error (if not, send response back)
	select {
	case err := <-errorChan:
		return err
	default:
		err := utils.SendFileResponse(fileContext)
		return err
	}
}

// /pdf-docx
func HandleDocxToPdf(fileContext *utils.FileContext) error {
	if err := fileContext.Ctx.BodyParser(fileContext.Request); err != nil {
		log.Println("Could not parse request body: ", err)
		return errors.New("Invalid request body")
	}
	// generate a unique hash as prefix for files so we have no collision
	fileContext.FilenamePrefix = utils.GenerateHash()

	// convert the files in parallel using go routines
	errorChan := make(chan error, len(fileContext.Request.Files))
	var wg sync.WaitGroup

	for _, file := range fileContext.Request.Files {
		wg.Add(1)
		go func(file *utils.File) {
			defer wg.Done()
			err := utils.ConvertFile(file, fileContext)
			if err != nil {
				errorChan <- err
			}
		}(file)
	}
	wg.Wait()

	// return err if any of the file convertings threw an error (if not, send response back)
	select {
	case err := <-errorChan:
		return err
	default:
		err := utils.SendFileResponse(fileContext)
		return err
	}
}
