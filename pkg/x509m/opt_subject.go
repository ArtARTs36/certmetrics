package x509m

import (
	"crypto/x509"
	"github.com/artarts36/certmetrics/pkg/collector"
)

type subjectNameOption struct {
	namer func(cert *x509.Certificate) string
}

func WithSubjectNameOf(namer func(cert *x509.Certificate) string) InspectOption {
	return &subjectNameOption{
		namer: namer,
	}
}

func WithoutSubjectName() InspectOption {
	return WithSubjectNameOf(func(_ *x509.Certificate) string {
		return "<hidden>"
	})
}

func (o *subjectNameOption) apply(cert *x509.Certificate, storing *collector.Cert) {
	storing.Subject = o.namer(cert)
}
