package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(filePath string) (*SyncMeiliConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config SyncMeiliConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
