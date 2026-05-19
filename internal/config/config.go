package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	DataBaseFilePath string `json:"DataBaseFilePath"`
}

func NewConfiguration(pathConfigFile string) (*Configuration, error) {
	data, error := os.ReadFile(pathConfigFile)
	if error != nil {
		return nil, error
	}
	var config Configuration
	error = json.Unmarshal(data, &config)
	if error != nil {
		return nil, error
	}
	return &config, nil
}
