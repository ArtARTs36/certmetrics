package jwtm

import "github.com/artarts36/certmetrics"

var defaultInspector = NewInspector(certmetrics.DefaultCollector)

func InspectToken(token string, opts ...InspectOption) error {
	return defaultInspector.InspectToken(token, opts...)
}

func InspectClaims(claims map[string]interface{}, opts ...InspectOption) {
	defaultInspector.InspectClaims(claims, opts...)
}
