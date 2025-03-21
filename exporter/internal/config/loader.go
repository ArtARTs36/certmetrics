package config

import (
	"fmt"
	"os"
	"strings"
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

	if cfg.Scrape.Interval <= 0 {
		cfg.Scrape.Interval = defaultInterval
	}

	if err = validate(&cfg); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	if err = injectEnv(&cfg); err != nil {
		return nil, fmt.Errorf("inject enviornment variables: %w", err)
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	for i, pem := range cfg.Scrape.X509.PEMs {
		if pem.Path == "" {
			return fmt.Errorf("scrape.x509.pems.%d.path required", i)
		}

		if !pem.Opts.Subject.Valid() {
			return fmt.Errorf(
				"scrape.x509.pems.%d.opts.subject have invalid value. Expected: [%s]",
				i,
				PemSubjectNameOptValues(),
			)
		}
	}

	return nil
}

func injectEnv(cfg *Config) error {
	for i, pem := range cfg.Scrape.X509.PEMs {
		if strings.HasPrefix(pem.Path, "$") {
			varName := strings.TrimPrefix(pem.Path, "$")

			if strings.HasPrefix(varName, "{") {
				varName = strings.TrimPrefix(varName, "{")
				varName = strings.TrimSuffix(varName, "}")
			}

			val, ok := os.LookupEnv(varName)
			if ok {
				return fmt.Errorf(
					"scrape.x509.pems.%d.path: environment variable %q not found",
					i,
					varName,
				)
			}

			pem.Path = val
		}
	}

	return nil
}
