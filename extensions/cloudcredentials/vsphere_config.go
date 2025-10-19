package cloudcredentials

import "github.com/rancher/shepherd/pkg/config"

// The json/yaml config key for the azure cloud credential config
const VSphereCredentialConfigurationFileKey = "vmwarevsphereCredentials"

// VSphereCredentialConfig is configuration need to create an vsphere cloud credential
type VSphereCredentialConfig struct {
	Password    string `json:"password" yaml:"password"`
	Username    string `json:"username" yaml:"username"`
	Vcenter     string `json:"vcenter" yaml:"vcenter"`
	VcenterPort string `json:"vcenterPort" yaml:"vcenterPort"`
}

// LoadVSphereCredentialConfig loads the VSphere credential config from the cattle_config file
func LoadVSphereCredentialConfig() CloudCredential {
	var cloudCredential CloudCredential
	var vsphereCredentialConfig VSphereCredentialConfig

	config.LoadConfig(VSphereCredentialConfigurationFileKey, &vsphereCredentialConfig)
	cloudCredential.VSphereCredentialConfig = &vsphereCredentialConfig

	return cloudCredential
}
