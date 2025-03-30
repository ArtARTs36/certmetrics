package config

type Config struct {
	HTTP struct {
		Addr string `yaml:"addr" json:"addr"`
	} `yaml:"http" json:"http"`

	Scrape ScrapeConfig `yaml:"scrape" json:"scrape"`
}

type ScrapeConfig struct {
	Interval Duration `yaml:"interval" json:"interval"`

	X509 struct {
		// Paths to .pem
		Paths []X509File `yaml:"paths" json:"paths"`
	} `yaml:"x509" json:"x509"`
}

type X509File struct {
	Path string `yaml:"path" json:"path"`
	ID   string `yaml:"id" json:"id"`
	Opts struct {
		Subject PemSubjectNameOpt `yaml:"subject" json:"subject"` // default: cn
	} `yaml:"opts" json:"opts"`
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
