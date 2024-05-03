package permutations

import (
	"encoding/json"
	"fmt"
)

// Relationship structs are used to connect related key/values together so that they can be permuted in sync
type Relationship struct {
	ParentValue       any      `json:"parentValue" yaml:"parentValue"`
	ChildKeyPath      []string `json:"childKeyPath" yaml:"childkeyPath"`
	ChildKeyPathValue any      `json:"childKeyPathValue" yaml:"childkeyPathValue"`
}

type PermutationObject struct {
	KeyPath                   []string       `json:"keyPath" yaml:"keyPath"`
	KeyPathValues             []any          `json:"keyPathValue" yaml:"keyPath"`
	KeyPathValueRelationships []Relationship `json:"keyPathValueRelationships" yaml:"KeyPathValueRelationships"`
}

func Permute(permutations []PermutationObject, baseConfig map[string]any) (permutedConfigs []map[string]any, permuteName string, err error) {
	var configs []map[string]any
	for _, keyPathValue := range permutations[0].KeyPathValues {
		marshaledConfig, err := json.Marshal(baseConfig)
		if err != nil {
			return nil, "", err
		}

		unmarshaledConfig := make(map[string]any)
		json.Unmarshal(marshaledConfig, &unmarshaledConfig)

		permutedConfig, err := ReplaceValue(permutations[0].KeyPath, keyPathValue, unmarshaledConfig)
		if err != nil {
			return nil, "", err
		}

		for _, relationship := range permutations[0].KeyPathValueRelationships {
			if relationship.ParentValue == keyPathValue {
				permutedConfig, err = ReplaceValue(relationship.ChildKeyPath, relationship.ChildKeyPathValue, permutedConfig)
				if err != nil {
					return nil, "", err
				}
			}
		}

		configs = append(configs, permutedConfig)
	}

	var finalConfigs []map[string]any
	if len(permutations) == 1 {
		return configs, "", nil
	} else {
		for _, config := range configs {
			permutedConfigs, _, err := Permute(permutations[1:], config)
			if err != nil {
				return nil, "", err
			}

			finalConfigs = append(finalConfigs, permutedConfigs...)
		}
	}

	return finalConfigs, "", err
}

func ReplaceValue(keyPath []string, replaceVal any, searchMap map[string]any) (map[string]any, error) {
	if len(keyPath) == 1 {
		searchMap[keyPath[0]] = replaceVal

		return searchMap, nil
	} else {
		var err error

		searchMap[keyPath[0]], err = ReplaceValue(keyPath[1:], replaceVal, searchMap[keyPath[0]].(map[string]any))
		if err != nil {
			return nil, err
		}
	}

	return searchMap, nil
}

func main() {
	etcd := map[string]any{
		"disableSnapshot":      false,
		"snapshotScheduleCron": "0 */5 * * *",
		"snapshotRetain":       3,
	}

	provisioningInput := map[string]any{
		"rke2KubernetesVersion": "v1.27.10+rke2r1",
		"cni":                   "calico",
		"providers":             "aws",
		"nodeProviders":         "defaultNodeProvider",
		"hardened":              false,
		"psact":                 "",
		"clusterSSHTests":       []string{"CheckCPU", "NodeReboot", "AuditLog"},
		"etcd":                  etcd,
	}

	config := map[string]any{
		"provisioningInput": provisioningInput,
	}

	awsRelationship0 := Relationship{
		ParentValue:       "aws",
		ChildKeyPath:      []string{"provisioningInput", "nodeProviders"},
		ChildKeyPathValue: "ec2",
	}

	awsRelationship1 := Relationship{
		ParentValue:       "aws",
		ChildKeyPath:      []string{"provisioningInput", "cloudProvider"},
		ChildKeyPathValue: "aws_cloud_provider",
	}

	azureRelationship0 := Relationship{
		ParentValue:       "azure",
		ChildKeyPath:      []string{"provisioningInput", "nodeProviders"},
		ChildKeyPathValue: "azure_node_provider",
	}

	testPermutationObject0 := PermutationObject{
		KeyPath:                   []string{"provisioningInput", "providers"},
		KeyPathValues:             []any{"aws", "azure", "vsphere"},
		KeyPathValueRelationships: []Relationship{awsRelationship0, awsRelationship1, azureRelationship0},
	}

	testPermutationObject1 := PermutationObject{
		KeyPath:                   []string{"provisioningInput", "rke2KubernetesVersion"},
		KeyPathValues:             []any{"v1.26.10+rke2r1", "v1.27.10+rke2r1", "v1.28.10+rke2r1"},
		KeyPathValueRelationships: []Relationship{},
	}

	permutedConfigs, _, err := Permute([]PermutationObject{testPermutationObject0, testPermutationObject1}, config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(permutedConfigs))
	for _, permutedConfig := range permutedConfigs {
		fmt.Println("---------------------------------------------------")
		indented, _ := json.MarshalIndent(permutedConfig, "", "  ")
		converted := string(indented)
		fmt.Println(converted)
	}
}
