package cloudcredentials

import "github.com/rancher/shepherd/pkg/config"

// The json/yaml config key for the aws cloud credential config
const AmazonCredentialConfigurationFileKey = "awsCredentials"

// AmazonCredentialConfig is configuration need to create an aws cloud credential
type AmazonCredentialConfig struct {
	AccessKey     string `json:"accessKey" yaml:"accessKey"`
	SecretKey     string `json:"secretKey" yaml:"secretKey"`
	DefaultRegion string `json:"defaultRegion,omitempty" yaml:"defaultRegion,omitempty"`
}

// LoadAmazonCredentialConfig loads the amazon credential config from the cattle_config file
func LoadAmazonCredentialConfig() CloudCredential {
	var cloudCredential CloudCredential
	var awsCredentialConfig AmazonCredentialConfig

	config.LoadConfig(AmazonCredentialConfigurationFileKey, &awsCredentialConfig)
	cloudCredential.AmazonCredentialConfig = &awsCredentialConfig

	return cloudCredential
}
