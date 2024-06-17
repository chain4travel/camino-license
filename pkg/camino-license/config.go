package caminolicense

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type PossibleHeader struct {
	Name   string `yaml:"name"`
	Header string `yaml:"header"`
}

type CustomHeader struct {
	Name   string   `yaml:"name"`
	Header string   `yaml:"header"`
	Paths  []string `yaml:"paths"`
}

type HeadersConfig struct {
	PossibleHeaders []PossibleHeader `yaml:"possible-headers"`
	CustomHeaders   []CustomHeader   `yaml:"custom-headers"`
}

func GetHeadersConfig(configPath string) (HeadersConfig, error) {

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return HeadersConfig{}, errors.Wrapf(err, "failed to read config file %s", configPath)
	}
	headersConfig := &HeadersConfig{}
	err = yaml.Unmarshal(yamlFile, headersConfig)
	if err != nil {
		return HeadersConfig{}, errors.Wrapf(err, "failed to read config file %s", configPath)
	}

	return *headersConfig, nil
}
