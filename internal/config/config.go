package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Target tek bir hedef site bilgisini temsil eder
type Target struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	Targets []Target `yaml:"targets"`
	Proxy   struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"proxy"`
	Timeout int `yaml:"timeout"`
}

func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Proxy.Host == "" {
		config.Proxy.Host = "127.0.0.1"
	}
	if config.Proxy.Port == "" {
		config.Proxy.Port = "9050"
	}
	if config.Timeout == 0 {
		config.Timeout = 60
	}

	return &config, nil
}
