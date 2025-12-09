package inits

import (
	"local-audio-lib/config"
	"os"

	"gopkg.in/yaml.v3"
)

func Config() (*config.Config, error) {
	// Read config file
	configFilePosition, exist := os.LookupEnv("CONFIG_FILE_PATH")
	if !exist {
		configFilePosition = "config.yml"
	}

	configFileBytes, err := os.ReadFile(configFilePosition)
	if err != nil {
		return nil, err
	}

	var cfg config.Config

	err = yaml.Unmarshal(configFileBytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
