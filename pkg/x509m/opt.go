package x509m

import (
	"crypto/x509"
	"github.com/artarts36/certmetrics/pkg/collector"
)

type InspectOption interface {
	apply(cert *x509.Certificate, storing *collector.Cert)
}
