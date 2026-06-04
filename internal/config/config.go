package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	DataBaseFilePath string `json:"DataBaseFilePath"`
}

func NewConfiguration(pathConfigFile string) (*Configuration, error) {
	data, err := os.ReadFile(pathConfigFile)
	if err != nil {
		return nil, err
	}
	var config Configuration
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
