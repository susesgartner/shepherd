package vsphere

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/shepherd/clients/rancher"
	v1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/cloudcredentials"
	"github.com/rancher/shepherd/extensions/defaults"
	"github.com/rancher/shepherd/extensions/defaults/stevetypes"
	"github.com/rancher/shepherd/extensions/steve"
	"github.com/rancher/shepherd/pkg/config"
	"github.com/rancher/shepherd/pkg/namegenerator"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	vsphereProvider     = "vsphere"
	credentialNamespace = "cattle-global-data"
)

// CreateVsphereCloudCredentials is a helper function that takes the rancher Client as a parameter and creates
// an AWS cloud credential, and returns the CloudCredential response
/*
func CreateVsphereCloudCredentials(rancherClient *rancher.Client) (*cloudcredentials.CloudCredential, error) {
	var vmwarevsphereCredentialConfig cloudcredentials.VmwarevsphereCredentialConfig
	config.LoadConfig(cloudcredentials.VmwarevsphereCredentialConfigurationFileKey, &vmwarevsphereCredentialConfig)

	cloudCredential := cloudcredentials.CloudCredential{
		Name:                vmwarevsphereCloudCredNameBase,
		VmwareVsphereConfig: &vmwarevsphereCredentialConfig,
	}

	resp := &cloudcredentials.CloudCredential{}
	err := rancherClient.Management.APIBaseClient.Ops.DoCreate(management.CloudCredentialType, cloudCredential, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}*/

func CreateVsphereCloudCredentials(client *rancher.Client, credentials cloudcredentials.CloudCredential) (*v1.SteveAPIObject, error) {
	users, err := client.Management.User.ListAll(&types.ListOpts{})
	if err != nil {
		return nil, err
	}

	secretName := namegenerator.AppendRandomString(vsphereProvider)
	spec := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "cc-",
			Namespace:    credentialNamespace,
			Annotations: map[string]string{
				"field.cattle.io/name":      secretName,
				"field.cattle.io/creatorId": users.Data[0].ID,
			},
		},
		Data: map[string][]byte{
			"vmwarevspherecredentialConfig-password":    []byte(credentials.VmwareVsphereConfig.Password),
			"vmwarevspherecredentialConfig-username":    []byte(credentials.VmwareVsphereConfig.Username),
			"vmwarevspherecredentialConfig-vcenter":     []byte(credentials.VmwareVsphereConfig.Vcenter),
			"vmwarevspherecredentialConfig-vcenterPort": []byte(credentials.VmwareVsphereConfig.VcenterPort),
		},
		Type: corev1.SecretTypeOpaque,
	}

	vSphereCloudCredentials, err := steve.CreateAndWaitForResource(client, stevetypes.Secret, spec, true, defaults.FiveSecondTimeout, defaults.FiveMinuteTimeout)
	if err != nil {
		return nil, err
	}

	return vSphereCloudCredentials, nil
}

// GetVspherePassword is a helper to get the password from the cloud credential object as a string
func GetVspherePassword() string {
	var vmwarevsphereCredentialConfig cloudcredentials.VmwarevsphereCredentialConfig

	config.LoadConfig(cloudcredentials.VmwarevsphereCredentialConfigurationFileKey, &vmwarevsphereCredentialConfig)

	return vmwarevsphereCredentialConfig.Password
}
