package cloudcredentials

import "github.com/rancher/shepherd/pkg/config"

// The json/yaml config key for the azure cloud credential config
const AzureCredentialConfigurationFileKey = "azureCredentials"

// AzureCredentialConfig is configuration need to create an azure cloud credential
type AzureCredentialConfig struct {
	ClientID       string `json:"clientId" yaml:"clientId"`
	ClientSecret   string `json:"clientSecret" yaml:"clientSecret"`
	SubscriptionID string `json:"subscriptionId" yaml:"subscriptionId"`
	Environment    string `json:"environment" yaml:"environment"`
}

// LoadAzureCredentialConfig loads the Azure credential config from the cattle_config file
func LoadAzureCredentialConfig() CloudCredential {
	var cloudCredential CloudCredential
	var azureCredentialConfig AzureCredentialConfig

	config.LoadConfig(AzureCredentialConfigurationFileKey, &azureCredentialConfig)
	cloudCredential.AzureCredentialConfig = &azureCredentialConfig

	return cloudCredential
}
