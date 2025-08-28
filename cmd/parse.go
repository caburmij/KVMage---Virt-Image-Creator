package cmd

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (map[string]*Options, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw struct {
		KVMage map[string]*Options `yaml:"kvmage"`
	}

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	return raw.KVMage, nil
}
