package config

const UploadthingVersion = "6.1.0"
const AppIdEnvKey = "UPLOADTHING_APP_ID"
const ApiKeyEnvKey = "UPLOADTHING_SECRET"
const CustomInfraUrlEnvKey = "CUSTOM_INFRA_URL"

type FileType string

const (
	Image FileType = "image"
	Video FileType = "video"
	Audio FileType = "audio"
	Text  FileType = "text"
	Pdf   FileType = "pdf"
	Blob  FileType = "blob"
)

type ContentDisposition string

const (
	Inline     ContentDisposition = "inline"
	Attachment ContentDisposition = "attachment"
)

type FileSize int

const (
	Bytes     FileSize = 1
	KiloBytes FileSize = 1024 * Bytes
	MegaBytes FileSize = 1024 * KiloBytes
	GigaBytes FileSize = 1024 * MegaBytes
)
