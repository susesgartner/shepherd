package cloudcredentials

import "github.com/rancher/shepherd/pkg/config"

// The json/yaml config key for the google cloud credential config
const GoogleCredentialConfigurationFileKey = "googleCredentials"

// GoogleCredentialConfig is configuration need to create a google cloud credential
type GoogleCredentialConfig struct {
	AuthEncodedJSON string `json:"authEncodedJson" yaml:"authEncodedJson"`
}

// LoadGoogleCredentialConfig loads the Google credential config from the cattle_config file
func LoadGoogleCredentialConfig() CloudCredential {
	var cloudCredential CloudCredential
	var googleCredentialConfig GoogleCredentialConfig

	config.LoadConfig(GoogleCredentialConfigurationFileKey, &googleCredentialConfig)
	cloudCredential.GoogleCredentialConfig = &googleCredentialConfig

	return cloudCredential
}
