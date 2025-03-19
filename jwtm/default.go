package jwtm

import "github.com/artarts36/certmetrics"

var defaultInspector = NewInspector(certmetrics.DefaultCollector)

func InspectToken(token string) error {
	return defaultInspector.InspectToken(token)
}

func InspectNamedToken(subjectName func(claims map[string]interface{}) string, token string) error {
	return defaultInspector.InspectNamedToken(subjectName, token)
}
