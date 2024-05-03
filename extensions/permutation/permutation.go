package permutation

import (
	"encoding/json"
	"os"

	"sigs.k8s.io/yaml"
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

func LoadConfigFile(filePath string) (map[string]any, error) {
	allString, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var all map[string]any
	err = yaml.Unmarshal(allString, &all)
	if err != nil {
		panic(err)
	}

	return all, nil
}
