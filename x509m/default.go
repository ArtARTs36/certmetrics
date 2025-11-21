package x509m

import "github.com/artarts36/certmetrics"

var defaultInspector = NewInspector(certmetrics.DefaultCollector)

func Inspect(certBytes []byte, opts ...InspectOption) error {
	return defaultInspector.Inspect(certBytes, opts...)
}

func InspectPEMs(pemCerts []byte, opts ...InspectOption) error {
	return defaultInspector.InspectPEMs(pemCerts, opts...)
}
