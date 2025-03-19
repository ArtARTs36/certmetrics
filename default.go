package certmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DefaultCollector = NewPrometheusCollector("go")
)

func Register() error {
	return prometheus.DefaultRegisterer.Register(DefaultCollector)
}
