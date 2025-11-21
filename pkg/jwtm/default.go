package jwtm

import (
	"github.com/artarts36/certmetrics/pkg/collector"
)

var defaultInspector = NewInspector(collector.DefaultCollector)

func InspectToken(token string, opts ...InspectOption) error {
	return defaultInspector.InspectToken(token, opts...)
}

func InspectClaims(claims map[string]interface{}, opts ...InspectOption) {
	defaultInspector.InspectClaims(claims, opts...)
}
