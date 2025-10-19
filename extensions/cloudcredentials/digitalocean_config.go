package cloudcredentials

import "github.com/rancher/shepherd/pkg/config"

// The json/yaml config key for the digital ocean cloud credential config
const DigitalOceanCredentialConfigurationFileKey = "digitalOceanCredentials"

// DigitalOceanCredentialConfig is configuration need to create a digital ocean cloud credential
type DigitalOceanCredentialConfig struct {
	AccessToken string `json:"accessToken" yaml:"accessToken"`
}

// LoadDigitalOceanCredentialConfig loads the DigitalOcean credential config from the cattle_config file
func LoadDigitalOceanCredentialConfig() CloudCredential {
	var cloudCredential CloudCredential
	var digitalOceanCredentialConfig DigitalOceanCredentialConfig

	config.LoadConfig(DigitalOceanCredentialConfigurationFileKey, &digitalOceanCredentialConfig)
	cloudCredential.DigitalOceanCredentialConfig = &digitalOceanCredentialConfig

	return cloudCredential
}
