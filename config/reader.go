package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Read() (Settings, error) {
	file, err := os.Open(".config/test.yaml")
	if err != nil {
		return Settings{}, fmt.Errorf("can't open file: %v", err)
	}

	settings := Settings{}
	err = yaml.NewDecoder(file).Decode(&settings)
	if err != nil {
		return Settings{}, fmt.Errorf("can't read settings: %v", err)
	}

	return settings, nil
}
