package x509m

import (
	"crypto/x509"

	"github.com/artarts36/certmetrics"
)

type idInspectOption struct {
	id string
}

func WithID(id string) InspectOption {
	return &idInspectOption{
		id: id,
	}
}

func (o *idInspectOption) apply(_ *x509.Certificate, storing *certmetrics.Cert) {
	storing.ID = o.id
}
