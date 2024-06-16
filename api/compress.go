package api

import (
	"errors"
	"fmt"
	"mime/multipart"
	"sync"

	"github.com/QuodEstDubitandum/iUtils/utils"
)

func HandleImgCompression(fileContext *utils.FileContext) error {
	form, err := fileContext.Ctx.MultipartForm()
	if err == nil {
		fileContext.Files = form.File["files"]
		fileContext.InputFormat = form.Value["input_format"][0]

		// if the inputformat and outputformat sent by request dont match allowed ones
		if _, ok := fileContext.AllowedInputFormats[fileContext.InputFormat]; !ok {
			fmt.Println("The allowed input seem to have been tampered with")
			return errors.New("The allowed input format seem to have been tampered with.")
		}

		// generate a unique hash as prefix for files so we have no collision
		fileContext.FilenamePrefix = utils.GenerateHash()

		// convert the files in parallel using go routines
		errorChan := make(chan error, len(fileContext.Files))
		var wg sync.WaitGroup

		for _, file := range fileContext.Files {
			wg.Add(1)
			go func(file *multipart.FileHeader) {
				defer wg.Done()
				err := utils.CompressFile(file, fileContext)
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
	return errors.New("The allowed content type seems to have been tampered with.")
}
