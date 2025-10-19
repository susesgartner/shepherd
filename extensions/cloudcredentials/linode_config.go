package cloudcredentials

import "github.com/rancher/shepherd/pkg/config"

// The json/yaml config key for the linode cloud credential config
const LinodeCredentialConfigurationFileKey = "linodeCredentials"

// LinodeCredentialConfig is configuration need to create a linode cloud credential
type LinodeCredentialConfig struct {
	Token string `json:"token" yaml:"token"`
}

// LoadLinodeCredentialConfig loads the Linode credential config from the cattle_config file
func LoadLinodeCredentialConfig() CloudCredential {
	var cloudCredential CloudCredential
	var linodeCredentialConfig LinodeCredentialConfig

	config.LoadConfig(LinodeCredentialConfigurationFileKey, &linodeCredentialConfig)
	cloudCredential.LinodeCredentialConfig = &linodeCredentialConfig

	return cloudCredential
}
