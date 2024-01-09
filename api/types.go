package api

import (
	"github.com/anuragkumar19/uploadthing-go/config"
)

type UploadFileMeta struct {
	Name string          `json:"name"`
	Type config.FileType `json:"type"`
	Size config.FileSize `json:"size"`
}

type UploadthingFile struct {
	PresignedUrls []string        `json:"presignedUrls"`
	Key           string          `json:"key"`
	FileUrl       string          `json:"fileUrl"`
	FileType      config.FileType `json:"fileType"`
	UploadId      string          `json:"uploadId"`
	ChunkSize     config.FileSize `json:"chunkSize"`
	ChunkCount    int             `json:"chunkCount"`
}

type UploadthingPresignedURLResponse struct {
	Data []UploadthingFile `json:"data"`
}

type UploadthingErrorResponse struct {
	Error string `json:"error"`
}

type ListFilesOptions struct {
	Limit  int
	Offset int
}

type RenameFile struct {
	FileKey string
	NewName string
}

type UploadthingError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Message string
	} `json:"data"`
}

func (err UploadthingError) Error() string {
	return err.Message
}

type EtagReturn struct {
	Etag       string
	Error      *UploadthingError
	PartNumber int
}

type UploadStatus struct {
	UploadthingFile
	Success bool
	Error   []UploadthingError
}
