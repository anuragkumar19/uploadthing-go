package api

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"

	"github.com/anuragkumar19/uploadthing-go/config"
	"github.com/gabriel-vasile/mimetype"
)

func (ut *UploadthingApi) UploadFiles(paths []string) ([]UploadStatus, error) {
	filesMeta := []UploadFileMeta{}

	for _, path := range paths {
		mtype, err := mimetype.DetectFile(path)

		if err != nil {
			return nil, err
		}

		metaData, err := os.Stat(path)

		if err != nil {
			return nil, err
		}

		filesMeta = append(filesMeta, UploadFileMeta{
			Name: metaData.Name(),
			Size: config.FileSize(metaData.Size()),
			Type: config.FileType(mtype.String()),
		})
	}

	uploadFiles, err := ut.getPresignedUrls(filesMeta)

	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	uploadStatusChan := make(chan UploadStatus)

	for i, uploadFile := range uploadFiles {
		wg.Add(1)
		go func(i int, uploadFile UploadthingFile) {
			defer wg.Done()

			wg2 := &sync.WaitGroup{}
			etagChan := make(chan EtagReturn)

			f, err := os.Open(paths[i])

			if err != nil {
				panic("file not found")
			}

			defer f.Close()

			reader := bufio.NewReader(f)
			chunkList := [][]byte{}
			buf := make([]byte, int(uploadFile.ChunkSize))

			for {
				n, err := reader.Read(buf)

				if err != nil {
					if err != io.EOF {
						log.Fatal(err)
					}
					break
				}
				chunkList = append(chunkList, buf[0:n])
			}

			for j, presignedUrl := range uploadFile.PresignedUrls {

				wg2.Add(1)
				go ut.uploadPart(wg2, etagChan, &UploadPartParams{
					Url:   presignedUrl,
					Key:   uploadFile.Key,
					Chunk: chunkList[j],
					//TODO:
					ContentType: string(uploadFile.FileType),
					//TODO:
					ContentDispositionType: "inline",
					FileName:               filesMeta[i].Name,
					MaxRetries:             10,
					PartNumber:             j + 1,
				}, 0)
			}

			go func() {
				defer close(etagChan)
				wg2.Wait()
			}()

			etags := []string{}
			errors := []UploadthingError{}

			for e := range etagChan {
				if e.Error != nil {
					// TODO: cancel all goroutine and return error to parent channel
					errors = append(errors, *e.Error)
				} else {
					etags = append(etags, e.Etag)
				}
			}

			if len(errors) == 0 {
				ut.completeMultipart(uploadFile.Key, uploadFile.UploadId, etags)
				success := ut.confirmUpload(uploadFile.Key, 20, 0)

				var maybeErrors = []UploadthingError{}

				if !success {
					maybeErrors = append(maybeErrors, UploadthingError{
						//TODO:
						Code:    "",
						Message: "Failed to verify upload",
					})
				}

				uploadStatusChan <- UploadStatus{
					UploadthingFile: uploadFile,
					Success:         success,
					Error:           maybeErrors,
				}
				return
			}

			uploadStatusChan <- UploadStatus{
				UploadthingFile: uploadFile,
				Success:         false,
				Error:           errors,
			}

		}(i, uploadFile)
	}

	go func() {
		defer close(uploadStatusChan)
		wg.Wait()
	}()

	results := []UploadStatus{}

	for r := range uploadStatusChan {
		results = append(results, r)
	}

	return results, nil
}

func (ut *UploadthingApi) UploadFilesFromURL(urls []string) {

}

func (ut *UploadthingApi) DeleteFiles(urls []string) {

}

func (ut *UploadthingApi) GetFileUrls(fileKeys []string) {

}

func (ut *UploadthingApi) ListFiles(option *ListFilesOptions) {

}

func (ut *UploadthingApi) RenameFiles(files []RenameFile) {

}
