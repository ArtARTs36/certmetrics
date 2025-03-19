package x509m

import (
	"crypto/x509"

	"github.com/artarts36/certmetrics"
)

type subjectNameOption struct {
	namer func(cert *x509.Certificate) string
}

func WithSubjectNameOf(namer func(cert *x509.Certificate) string) InspectOption {
	return &subjectNameOption{
		namer: namer,
	}
}

func (o *subjectNameOption) apply(cert *x509.Certificate, storing *certmetrics.Cert) {
	storing.Subject = o.namer(cert)
}
