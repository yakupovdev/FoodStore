package logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel  string `envconfig:"LEVEL"  required:"true"`
	LogFolder string `envconfig:"FOLDER" required:"true"`
}

func NewLoggerConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("log", &cfg); err != nil {
		return Config{}, fmt.Errorf("error parsing env vars: %w", err)
	}
	return cfg, nil
}

func NewConfigMust() Config {
	cfg, err := NewLoggerConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load logger config: %v", err))
	}
	return cfg
}
