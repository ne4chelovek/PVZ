package config

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func LoadConfig() (*Config, error) {
	env := getEnv()

	viper.SetConfigName(env)
	viper.SetConfigType("yml")

	absPath, err := filepath.Abs("configs")
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for configs: %w", err)
	}
	viper.AddConfigPath(absPath)

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AllowEmptyEnv(true)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if expStr, ok := viper.Get("jwt.expiration").(string); ok {
		duration, err := time.ParseDuration(expStr)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT expiration: %w", err)
		}
		cfg.JWT.Expiration = duration
	}

	log.Printf("Loaded config for environment: %s", env)
	return &cfg, nil
}

func getEnv() string {
	env := viper.GetString("ENV")
	if env == "" {
		env = "local"
	}
	return env
}
