package api

import (
	"net/http"
	"os"

	"github.com/anuragkumar19/uploadthing-go/config"
)

type UploadthingApiConfig struct {
	InfraURL string
	ApiKey   string
	TempDir  string
}

type UploadthingApi struct {
	infraURL string
	version  string
	apiKey   string
	tempDir  string
}

func (ut *UploadthingApi) getRequestUrl(path string) string {
	return ut.infraURL + path
}

func (ut *UploadthingApi) getDefaultHeaders() http.Header {
	return http.Header{
		"Content-Type":          {"application/json"},
		"x-uploadthing-api-key": {ut.apiKey},
		"x-uploadthing-version": {ut.version},
	}
}

func New() *UploadthingApi {
	return NewWithConfig(getDefaultConfig())
}

func NewWithConfig(conf *UploadthingApiConfig) *UploadthingApi {
	defaultConfig := getDefaultConfig()

	if conf.InfraURL != "" {
		defaultConfig.InfraURL = conf.InfraURL
	}

	if conf.TempDir != "" {
		defaultConfig.TempDir = conf.TempDir
	}

	if conf.ApiKey != "" {
		defaultConfig.ApiKey = conf.ApiKey
	}

	if defaultConfig.ApiKey == "" {
		panic("uploadthing: failed to load secret")
	}

	return &UploadthingApi{
		infraURL: defaultConfig.InfraURL,
		version:  config.UploadthingVersion,
		apiKey:   defaultConfig.ApiKey,
		tempDir:  defaultConfig.TempDir,
	}
}

func getDefaultConfig() *UploadthingApiConfig {
	infraUrl := os.Getenv(config.CustomInfraUrlEnvKey)

	if infraUrl == "" {
		infraUrl = "https://uploadthing.com"
	}

	apiKey := os.Getenv(config.ApiKeyEnvKey)

	conf := &UploadthingApiConfig{
		InfraURL: infraUrl,
		ApiKey:   apiKey,
		TempDir:  "tmp",
	}

	return conf
}
