package config

import (
	"time"
)

type Config struct {
	HTTP struct {
		Addr string `yaml:"addr"`
	} `yaml:"http"`

	Scrape struct {
		Interval time.Duration `yaml:"interval"`

		X509 struct {
			// Paths to .pem
			PEMs []PEMFile `yaml:"pems"`
		} `yaml:"x509"`
	} `yaml:"scrape"`
}

type PEMFile struct {
	Path string `yaml:"path"`
	ID   string `yaml:"id"`
}
