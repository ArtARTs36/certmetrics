package certmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DefaultCollector = NewPrometheusCollector("go")
	X509             = NewX509Inspector(DefaultCollector, nil)
	JWT              = NewJWTInspector(DefaultCollector)
)

func Register() error {
	return prometheus.DefaultRegisterer.Register(DefaultCollector)
}
