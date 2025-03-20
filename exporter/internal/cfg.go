package internal

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP struct {
		Addr string `yaml:"addr"`
	} `yaml:"http"`

	Inspect struct {
		Interval time.Duration `yaml:"interval"`

		X509 struct {
			// Paths to .pem
			PEMs []PEMFile `yaml:"pems"`
		} `yaml:"x509"`
	} `yaml:"inspect"`
}

type PEMFile struct {
	Path string `yaml:"path"`
	ID   string `yaml:"id"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	return &cfg, nil
}
