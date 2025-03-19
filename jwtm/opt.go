package jwtm

import (
	"github.com/artarts36/certmetrics"
)

type InspectOption interface {
	apply(claims map[string]interface{}, storing *certmetrics.Cert)
}
