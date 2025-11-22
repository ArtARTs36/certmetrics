package jwtm

import (
	"github.com/artarts36/certmetrics/pkg/collector"
)

type InspectOption interface {
	apply(claims map[string]interface{}, storing *collector.Cert)
}
