package x509m

import (
	"crypto/x509"

	"github.com/artarts36/certmetrics"
)

type InspectOption interface {
	apply(cert *x509.Certificate, storing *certmetrics.Cert)
}
