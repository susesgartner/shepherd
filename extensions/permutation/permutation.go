package permutation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Relationship structs are used to connect related key/values together so that they can be permuted in sync
type Relationship struct {
	ParentValue       any           `json:"parentValue" yaml:"parentValue"`
	ChildKeyPath      []string      `json:"childKeyPath" yaml:"childkeyPath"`
	ChildKeyPathValue any           `json:"childKeyPathValue" yaml:"childkeyPathValue"`
	ChildPermutations []Permutation `json:"childPermutations" yaml:"childPermutations"`
}

func CreateRelationship(parentValue any, childKeyPath []string, childKeyPathValue any, childPermutations []Permutation) Relationship {
	return Relationship{
		ParentValue:       parentValue,
		ChildKeyPath:      childKeyPath,
		ChildKeyPathValue: childKeyPathValue,
		ChildPermutations: childPermutations,
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

func Permute(permutations []Permutation, baseConfig map[string]any) ([]map[string]any, string, error) {
	var configs []map[string]any
	var err error
	if len(permutations) == 0 {
		return configs, "", err
	}

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

		subPermutations := false
		for _, relationship := range permutations[0].KeyPathValueRelationships {
			if relationship.ParentValue == keyPathValue {
				if len(relationship.ChildKeyPath) > 1 && relationship.ChildKeyPathValue != nil {
					permutedConfig, err = ReplaceValue(relationship.ChildKeyPath, relationship.ChildKeyPathValue, permutedConfig)
					if err != nil {
						return nil, "", err
					}
				}

				var relationshipPermutedConfigs []map[string]any
				if len(relationship.ChildPermutations) > 0 {
					subPermutations = true
					relationshipPermutedConfigs, _, err = Permute(relationship.ChildPermutations, permutedConfig)
					if err != nil {
						return nil, "", err
					}
				}
				configs = append(configs, relationshipPermutedConfigs...)
			}
		}

		if !subPermutations {
			configs = append(configs, permutedConfig)
		}
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
	if len(keyPath) <= 1 {
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

func GetKeyPathValue(keyPath []string, searchMap map[string]any) (any, error) {
	var err error
	var keypathvalues any
	if len(keyPath) == 1 {
		keypathvalues, ok := searchMap[keyPath[0]]
		if !ok {
			err = errors.New(fmt.Sprintf("expected key does not exist: %s", keyPath[0]))
		}
		return keypathvalues, err
	} else {
		if _, ok := searchMap[keyPath[0]].(map[string]any); ok {
			keypathvalues, err = GetKeyPathValue(keyPath[1:], searchMap[keyPath[0]].(map[string]any))
			if err != nil {
				return nil, err
			}
		} else if _, ok := searchMap[keyPath[0]].([]any); ok {
			for i := range searchMap[keyPath[0]].([]any) {
				keypathvalues, err = GetKeyPathValue(keyPath[1:], searchMap[keyPath[0]].([]any)[i].(map[string]any))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return keypathvalues, err
}
