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

func Init(path string) (*Config, error) {
	v := viper.New()

	fn := "config"
	ft := "json"

	v.SetConfigName(fn)
	v.SetConfigType(ft)
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("%s.%s not found, %w", fn, ft, err)
		} else {
			return nil, fmt.Errorf("read error, %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	return &cfg, nil
}
