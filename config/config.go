package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"os"
)

type Config struct {
	DBURL  string
	Host   string
	Port   int
	Secret string
}

func (c *Config) Validate() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Config - Validate: %w", err)
		}
	}()
	if c.DBURL == "" {
		return fmt.Errorf("missing DBURL field")
	} else if c.Host == "" {
		return fmt.Errorf("missing Host field")
	} else if c.Port == 0 {
		return fmt.Errorf("missing Port field")
	} else if c.Secret == "" {
		return fmt.Errorf("missing Secret field")
	}
	return
}

func Load() (*Config, error) {
	cfg := &Config{}

	cfg.DBURL = os.Getenv("DB_URL")
	cfg.Host = "0.0.0.0"
	cfg.Port = 8000
	cfg.Secret = os.Getenv("SECRET")
	if cfg.Secret == "" {
		cfg.Secret = encryptcookie.GenerateKey()
	}

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
