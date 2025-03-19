package certmetrics

import (
	"crypto/x509"
	"encoding/pem"
)

type X509IssuerNamer func(cert *x509.Certificate) string

type X509Inspector struct {
	collector Collector
	cfg       *X509Config
}

type X509Config struct {
	SubjectNamer func(cert *x509.Certificate) string
}

func setupX509Config(cfg *X509Config) *X509Config {
	if cfg == nil {
		cfg = &X509Config{}
	}

	if cfg.SubjectNamer == nil {
		cfg.SubjectNamer = func(cert *x509.Certificate) string {
			return cert.Subject.CommonName
		}
	}

	return cfg
}

func NewX509Inspector(collector Collector, cfg *X509Config) *X509Inspector {
	return &X509Inspector{
		collector: collector,
		cfg:       setupX509Config(cfg),
	}
}

func (i *X509Inspector) InspectPEMs(pemCerts []byte) {
	// based on: /go/src/crypto/x509/cert_pool.go:207

	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		certBytes := block.Bytes
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			continue
		}

		i.collector.StoreCert(i.cert(cert))
	}
}

func (i *X509Inspector) cert(cert *x509.Certificate) *Cert {
	return &Cert{
		Type:      "x509",
		Subject:   i.cfg.SubjectNamer(cert),
		StartedAt: cert.NotBefore,
		ExpiredAt: cert.NotAfter,
	}
}
