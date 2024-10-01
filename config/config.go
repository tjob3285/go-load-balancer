package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	URLs           []string `json:"urls"`
	Port           string   `json:"port"`
	Algorithm      string   `json:"algorithm"`
	HealthInterval string   `json:"healthInterval"`
}

func LoadConfig(file string) (Config, error) {
	var config Config

	data, err := os.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
