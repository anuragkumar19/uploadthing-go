package uploadthing

import "os"

type Uploadthing struct {
	appID  string
	secret string
}

type UploadthingApiConfig struct {
	InfraURL string
	ApiKey   string
	TempDir  string
}

func New() *Uploadthing {
	appIdEnvKey := "UPLOADTHING_APP_ID"
	secretEnvKey := "UPLOADTHING_APP_ID"

	appId := os.Getenv(appIdEnvKey)
	secret := os.Getenv(secretEnvKey)

	if appId == "" {
		panic("uploadthing: failed to load app id")
	}

	if secret == "" {
		panic("uploadthing: failed to load secret")
	}

	ut := Uploadthing{
		appID:  appId,
		secret: secret,
	}

	return &ut
}

func NewWithConfig(config *UploadthingConfig) *Uploadthing {

}
