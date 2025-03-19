package x509m

import "github.com/artarts36/certmetrics"

var defaultInspector = NewInspector(certmetrics.DefaultCollector, nil)

func InspectPEMs(pemCerts []byte) {
	defaultInspector.InspectPEMs(pemCerts)
}
