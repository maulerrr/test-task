package config

import (
	"fmt"
	"os"

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

	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Printf("Configuration file %s does not exist, falling back to environment variables\n", file)
	} else {
		fmt.Println("Loading configuration from file:", file)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		fmt.Println("Configuration file read successfully")
	}

	fmt.Println("Loading configuration from environment variables")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	fmt.Println("Configuration unmarshalled successfully")

	// fallback
	if config.Port == "" {
		config.Port = os.Getenv("PORT")
	}
	if config.JWTSecretKey == "" {
		config.JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	}
	if config.DBSource == "" {
		config.DBSource = os.Getenv("DB_SOURCE")
	}

	return &config, nil
}
