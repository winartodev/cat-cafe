package config

import (
	"fmt"
	"github.com/winartodev/cat-cafe/pkg/helper"
	"log"
	"os"
	"strconv"
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
	if err != nil {
		log.Printf("Warning: Config file not found at %s, using environment variables", developmentConfigPath)
	}

	OverrideConfig(&cfg)
	if cfg.App.Port == 0 {
		return nil, fmt.Errorf("config is empty: YAML missing and ENV vars not set")
	}

	return &cfg, nil
}

func OverrideConfig(cfg *Config) {
	if name := os.Getenv("APP_NAME"); name != "" {
		cfg.App.Name = name
	}
	if port := os.Getenv("APP_PORT"); port != "" {
		if p, err := strconv.ParseInt(port, 10, 32); err == nil {
			cfg.App.Port = int32(p)
		}
	}

	if driver := os.Getenv("DB_DRIVER"); driver != "" {
		cfg.Database.Driver = driver
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		cfg.Database.Port = port
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		cfg.Database.Name = name
	}
	if username := os.Getenv("DB_USERNAME"); username != "" {
		cfg.Database.Username = username
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		cfg.Database.SSLMode = sslMode
	}

	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		cfg.Redis.Addr = redisAddr
	}

	if db := os.Getenv("REDIS_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			cfg.Redis.DB = d
		}
	}

	if secretKey := os.Getenv("JWT_SECRET_KEY"); secretKey != "" {
		cfg.JWT.SecretKey = secretKey
	}
	if duration := os.Getenv("JWT_TOKEN_DURATION"); duration != "" {
		if d, err := strconv.ParseInt(duration, 10, 64); err == nil {
			cfg.JWT.TokenDuration = d
		}
	}
}
