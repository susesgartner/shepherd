package google

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/shepherd/clients/rancher"
	v1 "github.com/rancher/shepherd/clients/rancher/v1"
	"github.com/rancher/shepherd/extensions/cloudcredentials"
	"github.com/rancher/shepherd/extensions/defaults"
	"github.com/rancher/shepherd/extensions/defaults/stevetypes"
	"github.com/rancher/shepherd/extensions/steve"
	"github.com/rancher/shepherd/pkg/namegenerator"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	googleProvider      = "gcp"
	credentialNamespace = "cattle-global-data"
)

// CreateGoogleCloudCredentials is a helper function that creates V1 cloud credentials and waits for them to become active.
func CreateGoogleCloudCredentials(client *rancher.Client, credentials cloudcredentials.CloudCredential) (*v1.SteveAPIObject, error) {
	users, err := client.Management.User.ListAll(&types.ListOpts{})
	if err != nil {
		return nil, err
	}

	secretName := namegenerator.AppendRandomString(googleProvider)
	spec := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "cc-",
			Namespace:    credentialNamespace,
			Annotations: map[string]string{
				"field.cattle.io/name":          secretName,
				"provisioning.cattle.io/driver": googleProvider,
				"field.cattle.io/creatorId":     users.Data[0].ID,
			},
		},
		Data: map[string][]byte{
			"googlecredentialConfig-authEncodedJson": []byte(credentials.GoogleCredentialConfig.AuthEncodedJSON),
		},
		Type: corev1.SecretTypeOpaque,
	}

	googleCloudCredentials, err := steve.CreateAndWaitForResource(client, stevetypes.Secret, spec, true, defaults.FiveSecondTimeout, defaults.FiveMinuteTimeout)
	if err != nil {
		return nil, err
	}

	return googleCloudCredentials, nil
}
