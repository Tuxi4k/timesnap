package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Path    string `mapstructure:"path"`
	Migrate bool   `mapstructure:"migrate"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

func Load() (*Config, error) {
	var cfg Config

	appMode := os.Getenv("APP_MODE")
	if appMode != "prod" {
		appMode = "dev"
	}

	configFile := fmt.Sprintf("config.%s", appMode)

	v := viper.New()
	v.SetConfigName(configFile)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.SetDefault("server.port", 8000)
	v.SetDefault("database.path", "db.sqlite3")
	v.SetDefault("database.migrate", true)

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
