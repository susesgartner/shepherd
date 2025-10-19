package cloudcredentials

import "github.com/rancher/shepherd/pkg/config"

// The json/yaml config key for the alibaba cloud credential config
const AlibabaCredentialConfigurationFileKey = "alibabaCredentials"

// AlibabaCredentialConfig is configuration need to create an alibaba cloud credential
type AlibabaCredentialConfig struct {
	AccessKeyId     string `json:"accessKeyId" yaml:"accessKeyId"`
	SecretAccessKey string `json:"accessKeySecret" yaml:"accessKeySecret"`
}

// LoadAlibabaCredentialConfig loads the alibaba credential config from the cattle_config file
func LoadAlibabaCredentialConfig() CloudCredential {
	var cloudCredential CloudCredential
	var alibabaCredentialConfig AlibabaCredentialConfig

	config.LoadConfig(AlibabaCredentialConfigurationFileKey, &alibabaCredentialConfig)
	cloudCredential.AlibabaCredentialConfig = &alibabaCredentialConfig

	return cloudCredential
}
