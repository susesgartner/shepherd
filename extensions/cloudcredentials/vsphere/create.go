package vsphere

import (
	"github.com/rancher/shepherd/clients/rancher"
	v1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/cloudcredentials"
	"github.com/rancher/shepherd/extensions/defaults"
	"github.com/rancher/shepherd/extensions/defaults/namespaces"
	"github.com/rancher/shepherd/extensions/defaults/providers"
	"github.com/rancher/shepherd/extensions/defaults/stevestates"
	"github.com/rancher/shepherd/extensions/defaults/stevetypes"
	"github.com/rancher/shepherd/extensions/steve"
	"github.com/rancher/shepherd/pkg/config"
	"github.com/rancher/shepherd/pkg/namegenerator"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	localCluster = "local"
)

// CreateVsphereCloudCredentials is a helper function that creates V1 cloud credentials and waits for them to become active.
func CreateVsphereCloudCredentials(client *rancher.Client, credentials cloudcredentials.CloudCredential) (*v1.SteveAPIObject, error) {
	secretName := namegenerator.AppendRandomString(providers.Vsphere)
	spec := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cloudcredentials.GeneratedName,
			Namespace:    namespaces.CattleData,
			Annotations: map[string]string{
				"field.cattle.io/name":      secretName,
				"field.cattle.io/creatorId": client.UserID,
			},
		},
		Data: map[string][]byte{
			"vmwarevspherecredentialConfig-password":    []byte(credentials.VSphereCredentialConfig.Password),
			"vmwarevspherecredentialConfig-username":    []byte(credentials.VSphereCredentialConfig.Username),
			"vmwarevspherecredentialConfig-vcenter":     []byte(credentials.VSphereCredentialConfig.Vcenter),
			"vmwarevspherecredentialConfig-vcenterPort": []byte(credentials.VSphereCredentialConfig.VcenterPort),
		},
		Type: corev1.SecretTypeOpaque,
	}

	vSphereCloudCredentials, err := steve.CreateAndWaitForResource(client, namespaces.FleetLocal+"/"+localCluster, stevetypes.Secret, spec, stevestates.Active, defaults.FiveSecondTimeout, defaults.FiveMinuteTimeout)
	if err != nil {
		return nil, err
	}

	return vSphereCloudCredentials, nil
}

// GetVspherePassword is a helper to get the password from the cloud credential object as a string
func GetVspherePassword() string {
	var vmwarevsphereCredentialConfig cloudcredentials.VSphereCredentialConfig

	config.LoadConfig(cloudcredentials.VSphereCredentialConfigurationFileKey, &vmwarevsphereCredentialConfig)

	return vmwarevsphereCredentialConfig.Password
}
