package jwtm

import (
	"github.com/artarts36/certmetrics/pkg/collector"
)

type idInspectOption struct {
	id string
}

func WithID(id string) InspectOption {
	return &idInspectOption{
		id: id,
	}
}

func (o *idInspectOption) apply(_ map[string]interface{}, storing *collector.Cert) {
	storing.ID = o.id
}
