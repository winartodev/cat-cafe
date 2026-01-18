package config

import (
	"github.com/winartodev/cat-cafe/pkg/helper"
)

const (
	developmentConfigPath = "./config.yaml"
)

type JWTConfig struct {
	SecretKey     string `yaml:"secretKey"`
	TokenDuration int64  `yaml:"tokenDuration"`
}

type Config struct {
	App struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		Port int32  `yaml:"port"`
	} `yaml:"app"`

	Database Database    `yaml:"database"`
	Redis    RedisConfig `yaml:"redis"`
	JWT      JWTConfig   `yaml:"jwt"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := helper.ReadYaml(developmentConfigPath, &cfg)
	return &cfg, err
}
