package x509m

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/artarts36/certmetrics"
)

type Inspector struct {
	collector certmetrics.Collector
}

// NewInspector creates new instance of Inspector. cfg not required.
func NewInspector(collector certmetrics.Collector) *Inspector {
	return &Inspector{
		collector: collector,
	}
}

func (i *Inspector) InspectPEMs(pemCerts []byte, opts ...InspectOption) {
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

		i.collector.StoreCert(i.cert(cert, opts))
	}
}

func (i *Inspector) cert(cert *x509.Certificate, opts []InspectOption) *certmetrics.Cert {
	storing := &certmetrics.Cert{
		Type:      "x509",
		Subject:   cert.Subject.CommonName,
		StartedAt: cert.NotBefore,
		ExpiredAt: cert.NotAfter,
	}

	for _, opt := range opts {
		opt.apply(cert, storing)
	}

	return storing
}
