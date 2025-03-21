package x509m

import "github.com/artarts36/certmetrics"

var defaultInspector = NewInspector(certmetrics.DefaultCollector)

func InspectPEMs(pemCerts []byte, opts ...InspectOption) error {
	return defaultInspector.InspectPEMs(pemCerts, opts...)
}
