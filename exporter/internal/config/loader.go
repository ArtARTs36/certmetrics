package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultInterval = 24 * time.Hour

func Parse(content []byte) (*Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	injectEnv(&cfg)
	defaults(&cfg)

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return &cfg, nil
}

func defaults(cfg *Config) {
	if cfg.HTTP.Addr == "" {
		cfg.HTTP.Addr = ":8010"
	}

	if cfg.Scrape.Interval <= 0 {
		cfg.Scrape.Interval = defaultInterval
	}
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

func injectEnv(cfg *Config) {
	for _, pem := range cfg.Scrape.X509.PEMs {
		pem.Path = interpolateEnv(pem.Path)
	}

	cfg.HTTP.Addr = interpolateEnv(cfg.HTTP.Addr)
}

func interpolateEnv(str string) string {
	if !strings.HasPrefix(str, "$") {
		return str
	}

	varName := strings.TrimPrefix(str, "$")

	if strings.HasPrefix(varName, "{") {
		varName = strings.TrimPrefix(varName, "{")
		varName = strings.TrimSuffix(varName, "}")
	}

	val, ok := os.LookupEnv(varName)
	if !ok {
		slog.With(slog.String("var_name", varName)).Debug("[config] environment variable not found")

		return ""
	}

	return val
}
