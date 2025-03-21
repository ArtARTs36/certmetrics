package certmetrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusCollector struct {
	certInfo       *prometheus.GaugeVec
	certExpiryDays *prometheus.GaugeVec
	reporter       string
}

func NewPrometheusCollector(namespace string) *PrometheusCollector {
	return &PrometheusCollector{
		reporter: "library",
		certExpiryDays: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "cert_expiry_days",
			Namespace: namespace,
			Help:      "Days until certificate expiry",
		}, []string{
			"id",
			"subject",
		}),
		certInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "cert_info",
			Namespace: namespace,
			Help:      "Cert info",
		}, []string{
			"id",
			"type",
			"subject",
			"started_at",
			"expired_at",
			"reporter",
		}),
	}
}

func (c *PrometheusCollector) As(reporterName string) {
	c.reporter = reporterName
}

func (c *PrometheusCollector) Describe(desc chan<- *prometheus.Desc) {
	c.certInfo.Describe(desc)
	c.certExpiryDays.Describe(desc)
}

func (c *PrometheusCollector) Collect(metric chan<- prometheus.Metric) {
	c.certInfo.Collect(metric)
	c.certExpiryDays.Collect(metric)
}

func (c *PrometheusCollector) StoreCert(cert *Cert) {
	startedAt := "<unknown>"
	expiredAt := "<unknown>"

	if !cert.StartedAt.IsZero() {
		startedAt = cert.StartedAt.Format(time.DateTime)
	}
	if !cert.ExpiredAt.IsZero() {
		expiredAt = cert.ExpiredAt.Format(time.DateTime)

		c.storeExpiryDays(cert)
	}

	c.certInfo.WithLabelValues(
		cert.ID,
		cert.Type,
		cert.Subject,
		startedAt,
		expiredAt,
		c.reporter,
	).Set(1)
}

func (c *PrometheusCollector) storeExpiryDays(cert *Cert) {
	const dayHour = 24

	days := int64(time.Until(cert.ExpiredAt).Hours() / dayHour)

	c.certExpiryDays.WithLabelValues(cert.ID, cert.Subject).Set(float64(days))
}
