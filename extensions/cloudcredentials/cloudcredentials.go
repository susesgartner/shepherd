package cloudcredentials

import (
	"fmt"

	"github.com/rancher/norman/types"
	"github.com/rancher/shepherd/extensions/provisioninginput"
	"github.com/rancher/shepherd/pkg/config"
)

// CloudCredential is the main struct needed to create a cloud credential depending on the outside cloud service provider
type CloudCredential struct {
	types.Resource
	Annotations                  map[string]string              `json:"annotations,omitempty"`
	Created                      string                         `json:"created,omitempty"`
	CreatorID                    string                         `json:"creatorId,omitempty"`
	Description                  string                         `json:"description,omitempty"`
	Labels                       map[string]string              `json:"labels,omitempty"`
	Name                         string                         `json:"name,omitempty"`
	Removed                      string                         `json:"removed,omitempty"`
	AmazonEC2CredentialConfig    *AmazonEC2CredentialConfig     `json:"amazonec2credentialConfig,omitempty"`
	AzureCredentialConfig        *AzureCredentialConfig         `json:"azurecredentialConfig,omitempty"`
	DigitalOceanCredentialConfig *DigitalOceanCredentialConfig  `json:"digitaloceancredentialConfig,omitempty"`
	LinodeCredentialConfig       *LinodeCredentialConfig        `json:"linodecredentialConfig,omitempty"`
	HarvesterCredentialConfig    *HarvesterCredentialConfig     `json:"harvestercredentialConfig,omitempty"`
	GoogleCredentialConfig       *GoogleCredentialConfig        `json:"googlecredentialConfig,omitempty"`
	VmwareVsphereConfig          *VmwarevsphereCredentialConfig `json:"vmwarevspherecredentialConfig,omitempty"`
	UUID                         string                         `json:"uuid,omitempty"`
}

func LoadCloudCredential(provider string) CloudCredential {
	var cloudCredential CloudCredential
	switch {

	case provider == provisioninginput.AWSProviderName.String():
		var awsCredentialConfig AmazonEC2CredentialConfig
		config.LoadConfig(AmazonEC2CredentialConfigurationFileKey, &awsCredentialConfig)
		cloudCredential.AmazonEC2CredentialConfig = &awsCredentialConfig
		return cloudCredential

	case provider == provisioninginput.AzureProviderName.String():
		var azureCredentialConfig AzureCredentialConfig
		config.LoadConfig(AzureCredentialConfigurationFileKey, &azureCredentialConfig)
		cloudCredential.AzureCredentialConfig = &azureCredentialConfig
		return cloudCredential

	case provider == provisioninginput.DOProviderName.String():
		var digitalOceanCredentialConfig DigitalOceanCredentialConfig
		config.LoadConfig(DigitalOceanCredentialConfigurationFileKey, &digitalOceanCredentialConfig)
		cloudCredential.DigitalOceanCredentialConfig = &digitalOceanCredentialConfig
		return cloudCredential

	case provider == provisioninginput.LinodeProviderName.String():
		var linodeCredentialConfig LinodeCredentialConfig
		config.LoadConfig(LinodeCredentialConfigurationFileKey, &linodeCredentialConfig)
		cloudCredential.LinodeCredentialConfig = &linodeCredentialConfig
		return cloudCredential

	case provider == provisioninginput.HarvesterProviderName.String():
		var harvesterCredentialConfig HarvesterCredentialConfig
		config.LoadConfig(HarvesterCredentialConfigurationFileKey, &harvesterCredentialConfig)
		cloudCredential.HarvesterCredentialConfig = &harvesterCredentialConfig
		return cloudCredential

	case provider == provisioninginput.VsphereProviderName.String():
		var vsphereCredentialConfig VmwarevsphereCredentialConfig
		config.LoadConfig(VmwarevsphereCredentialConfigurationFileKey, &vsphereCredentialConfig)
		cloudCredential.VmwareVsphereConfig = &vsphereCredentialConfig
		return cloudCredential

	case provider == provisioninginput.GoogleProviderName.String():
		var googleCredentialConfig GoogleCredentialConfig
		config.LoadConfig(GoogleCredentialConfigurationFileKey, &googleCredentialConfig)
		cloudCredential.GoogleCredentialConfig = &googleCredentialConfig
		return cloudCredential

	default:
		panic(fmt.Sprintf("Provider:%v not found", provider))
	}
}
