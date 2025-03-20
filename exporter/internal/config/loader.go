package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultInterval = 24 * time.Hour

func Load(path string) (*Config, error) {
	var cfg Config

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	if cfg.HTTP.Addr == "" {
		cfg.HTTP.Addr = ":8010"
	}

	if cfg.Scrape.Interval == 0 {
		cfg.Scrape.Interval = defaultInterval
	}

	return &cfg, nil
}
