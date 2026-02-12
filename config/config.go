package config

import (
	"fmt"
	"os"

	"go-simpler.org/env"
)

type Config struct {
	Port int `env:"PORT" default:"8080"`

	DB struct {
		File string `env:"FILE"`
	} `env:"DB_"`

	Sentry struct {
		APIKey           string `env:"API_KEY,required"`
		Endpoint         string `env:"ENDPOINT" default:"http://sentry.io/api/0/"`
		OrganizationSlug string `env:"ORGANIZATION_SLUG"`
	} `env:"SENTRY_"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Load(&cfg, nil); err != nil {
		fmt.Println("Usage:")
		env.Usage(&cfg, os.Stdout, nil)
		return &cfg, err
	}

	return &cfg, nil
}
