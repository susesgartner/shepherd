package permutation

import (
	"encoding/json"
	"errors"
	"maps"
	"os"

	"sigs.k8s.io/yaml"
)

// Relationship structs are used to connect related key/values together so that they can be permuted in sync
type Relationship struct {
	ParentValue       any      `json:"parentValue" yaml:"parentValue"`
	ChildKeyPath      []string `json:"childKeyPath" yaml:"childkeyPath"`
	ChildKeyPathValue any      `json:"childKeyPathValue" yaml:"childkeyPathValue"`
}

func CreateRelationship(parentValue any, childKeyPath []string, childKeyPathValue any) Relationship {
	return Relationship{
		ParentValue:       parentValue,
		ChildKeyPath:      childKeyPath,
		ChildKeyPathValue: childKeyPathValue,
	}
}

type Permutation struct {
	KeyPath                   []string       `json:"keyPath" yaml:"keyPath"`
	KeyPathValues             []any          `json:"keyPathValue" yaml:"keyPath"`
	KeyPathValueRelationships []Relationship `json:"keyPathValueRelationships" yaml:"KeyPathValueRelationships"`
}

func CreatePermutation(keyPath []string, keyPathValues []any, keyPathValueRelationships []Relationship) Permutation {
	return Permutation{
		KeyPath:                   keyPath,
		KeyPathValues:             keyPathValues,
		KeyPathValueRelationships: keyPathValueRelationships,
	}
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

		if _, ok := searchMap[keyPath[0]].(map[string]any); ok {
			searchMap[keyPath[0]], err = ReplaceValue(keyPath[1:], replaceVal, searchMap[keyPath[0]].(map[string]any))
			if err != nil {
				return nil, err
			}
		} else if _, ok := searchMap[keyPath[0]].([]any); ok {
			for i := range searchMap[keyPath[0]].([]any) {
				searchMap[keyPath[0]].([]any)[i], err = ReplaceValue(keyPath[1:], replaceVal, searchMap[keyPath[0]].([]any)[i].(map[string]any))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return searchMap, nil
}

func LoadConfigFromFile(filePath string) (map[string]any, error) {
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

func LoadConfigFromMap(key string, configMap map[string]any, config interface{}) {
	scoped := configMap[key]
	scopedString, err := yaml.Marshal(scoped)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(scopedString, config)
	if err != nil {
		panic(err)
	}
}

// TODO move this to another file as this function is modular
func LoadDefaults(filePath string, overrideConfig map[string]any) (map[string]any, error) {
	if filePath == "" {
		yaml.Unmarshal([]byte("{}"), filePath)
		err := errors.New("No default file found")
		return nil, err
	}

	config, err := LoadConfigFromFile(filePath)
	if err != nil {
		return nil, err
	}

	//Override the default values with any provided values
	maps.Copy(overrideConfig, config)

	return config, nil
}
