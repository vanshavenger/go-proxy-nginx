package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

// ParseYAMLConfig reads a YAML file and returns its contents as a string
func ParseYAMLConfig(filePath string) (string, error) {
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading YAML config file: %w", err)
	}

	var config RootConfig

	err = yaml.Unmarshal(fileContents, &config)

	if err != nil {
		return "", fmt.Errorf("error unmarshalling YAML config file: %w", err)
	}

	jsonBytes, err := json.Marshal(config)

	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}

	return string(jsonBytes), nil

}

// ValidateConfig validates the JSON configuration
func ValidateConfig(configJSON string) (*RootConfig, error) {
	var config RootConfig
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON config: %w", err)
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}
