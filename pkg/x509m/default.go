package x509m

import (
	"github.com/artarts36/certmetrics/pkg/collector"
)

var defaultInspector = NewInspector(collector.DefaultCollector)

func Inspect(certBytes []byte, opts ...InspectOption) error {
	return defaultInspector.Inspect(certBytes, opts...)
}

func InspectPEMs(pemCerts []byte, opts ...InspectOption) error {
	return defaultInspector.InspectPEMs(pemCerts, opts...)
}
