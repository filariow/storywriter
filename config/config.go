package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const DefaultConfigFile string = ".storywriter.yaml"

func ParseDefault() (*Config, error) {
	return Parse(DefaultConfigFile)
}

func Parse(configFile string) (*Config, error) {
	v, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer v.Close()

	c, err := io.ReadAll(v)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	if err := yaml.Unmarshal(c, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Config struct {
	Templates Templates              `yaml:"templates"`
	Output    Output                 `yaml:"output"`
	Defaults  map[string]interface{} `yaml:"defaults"`
}

type Templates struct {
	Folder string `yaml:"folder"`
}

type Output struct {
	Folder string                 `yaml:"folder"`
	Typed  map[string]interface{} `yaml:"typed"`
}
