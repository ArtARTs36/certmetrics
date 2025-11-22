package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DefaultCollector = NewPrometheusCollector("certmetrics")
)

func Register() error {
	return prometheus.DefaultRegisterer.Register(DefaultCollector)
}
