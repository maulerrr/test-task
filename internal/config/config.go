package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"PORT"`
	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`
	DBSource     string `mapstructure:"DB_SOURCE"`
}

func LoadConfig(file string) (*Config, error) {
	viper.SetConfigFile(file)
	viper.AutomaticEnv()

	fmt.Println("Loading configuration from file:", file)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	fmt.Println("Configuration file read successfully")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	fmt.Println("Configuration unmarshalled successfully")

	return &config, nil
}
