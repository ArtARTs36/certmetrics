package config

import (
	"time"
)

type Config struct {
	HTTP struct {
		Addr string `yaml:"addr"`
	} `yaml:"http"`

	Scrape ScrapeConfig `yaml:"scrape"`
}

type ScrapeConfig struct {
	Interval time.Duration `yaml:"interval"`

	X509 struct {
		// Paths to .pem
		PEMs []PEMFile `yaml:"pems"`
	} `yaml:"x509"`
}

type PEMFile struct {
	Path string `yaml:"path"`
	ID   string `yaml:"id"`
	Opts struct {
		Subject PemSubjectNameOpt `yaml:"subject"` // default: cn
	} `yaml:"opts"`
}

type PemSubjectNameOpt string

const (
	PemSubjectNameOptUnspecified PemSubjectNameOpt = ""
	PemSubjectNameOptCn          PemSubjectNameOpt = "cn"
	PemSubjectNameOptNone        PemSubjectNameOpt = "none"
)

func PemSubjectNameOptValues() []PemSubjectNameOpt {
	return []PemSubjectNameOpt{
		PemSubjectNameOptCn,
		PemSubjectNameOptNone,
	}
}

func (o PemSubjectNameOpt) Valid() bool {
	switch o {
	case PemSubjectNameOptUnspecified:
		return true
	case PemSubjectNameOptNone:
		return true
	case PemSubjectNameOptCn:
		return true
	}

	return false
}
