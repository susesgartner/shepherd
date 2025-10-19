package cloudcredentials

import "github.com/rancher/shepherd/pkg/config"

// The json/yaml config key for the harvester cloud credential config
const HarvesterCredentialConfigurationFileKey = "harvesterCredentials"

// HarvesterCredentialConfig is configuration need to create a harvester cloud credential
type HarvesterCredentialConfig struct {
	ClusterID         string `json:"clusterId" yaml:"clusterId"`
	ClusterType       string `json:"clusterType" yaml:"clusterType"`
	KubeconfigContent string `json:"kubeconfigContent" yaml:"kubeconfigContent"`
}

// LoadHarvesterCredentialConfig loads the Harvester credential config from the cattle_config file
func LoadHarvesterCredentialConfig() CloudCredential {
	var cloudCredential CloudCredential
	var harvesterCredentialConfig HarvesterCredentialConfig

	config.LoadConfig(HarvesterCredentialConfigurationFileKey, &harvesterCredentialConfig)
	cloudCredential.HarvesterCredentialConfig = &harvesterCredentialConfig

	return cloudCredential
}
