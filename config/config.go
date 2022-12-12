package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBURL string
	Host  string
	Port  int
}

func (c *Config) Validate() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Config - Validate: %w", err)
		}
	}()
	if c.DBURL == "" {
		return fmt.Errorf("missing DBURL field")
	}
	return
}

func Load() (*Config, error) {
	cfg := &Config{}

	cfg.DBURL = os.Getenv("DB_URL")
	cfg.Host = "0.0.0.0"
	cfg.Port = 8000

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
