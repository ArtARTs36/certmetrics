package certmetrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusCollector struct {
	certInfo *prometheus.GaugeVec
}

func NewPrometheusCollector(namespace string) *PrometheusCollector {
	return &PrometheusCollector{
		certInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "cert_info",
			Namespace: namespace,
			Help:      "Cert info",
		}, []string{
			"type",
			"subject",
			"started_at",
			"expired_at",
		}),
	}
}

func (c *PrometheusCollector) Describe(desc chan<- *prometheus.Desc) {
	c.certInfo.Describe(desc)
}

func (c *PrometheusCollector) Collect(metric chan<- prometheus.Metric) {
	c.certInfo.Collect(metric)
}

func (c *PrometheusCollector) StoreCert(cert *Cert) {
	c.certInfo.WithLabelValues(
		cert.Type,
		cert.Subject,
		cert.StartedAt.Format(time.DateTime),
		cert.ExpiredAt.Format(time.DateTime),
	).Set(1)
}
