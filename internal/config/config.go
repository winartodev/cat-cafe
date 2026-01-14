package config

import (
	"github.com/winartodev/cat-cafe/pkg/helper"
)

const (
	developmentConfigPath = "./config.yaml"
)

type Config struct {
	App struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		Port int32  `yaml:"port"`
	} `yaml:"app"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := helper.ReadYaml(developmentConfigPath, &cfg)
	return &cfg, err
}
