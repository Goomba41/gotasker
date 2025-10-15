package configuration

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
	SslMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

func Init() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("Not found: add config.json to app root directory")
		} else {
			return nil, fmt.Errorf("Read configuration file, %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Parse error: %w", err)
	}

	return &cfg, nil
}
