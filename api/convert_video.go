package api

import (
	"errors"
	"mime/multipart"
	"sync"

	"github.com/QuodEstDubitandum/iUtils/utils"
)

// /mp4-gif
func HandleMP4ToGIF(fileContext *utils.FileContext) error {
	form, err := fileContext.Ctx.MultipartForm()
	if err == nil {
		fileContext.Files = form.File["files"]
		fileContext.InputFormat = "video/mp4"
		fileContext.OutputFormat = "image/gif"

		// generate a unique hash as prefix for files so we have no collision
		fileContext.FilenamePrefix = utils.GenerateHash()

		// convert the files in parallel using go routines
		errorChan := make(chan error, len(fileContext.Files))
		var wg sync.WaitGroup

		for _, file := range fileContext.Files {
			wg.Add(1)
			go func(file *multipart.FileHeader) {
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
	return errors.New("The allowed content type seems to have been tampered with.")
}

// /mp4-mp3
func HandleMP4ToMP3(fileContext *utils.FileContext) error {
	form, err := fileContext.Ctx.MultipartForm()
	if err == nil {
		fileContext.Files = form.File["files"]
		fileContext.InputFormat = "video/mp4"
		fileContext.OutputFormat = "audio/mpeg"

		// generate a unique hash as prefix for files so we have no collision
		fileContext.FilenamePrefix = utils.GenerateHash()

		// convert the files in parallel using go routines
		errorChan := make(chan error, len(fileContext.Files))
		var wg sync.WaitGroup

		for _, file := range fileContext.Files {
			wg.Add(1)
			go func(file *multipart.FileHeader) {
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
	return errors.New("The allowed content type seems to have been tampered with.")
}
