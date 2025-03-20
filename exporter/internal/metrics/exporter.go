package metrics

import "github.com/prometheus/client_golang/prometheus"

type ExporterMetrics struct {
	scrapingsTotal *prometheus.CounterVec
}

func NewExporterMetrics(namespace string) *ExporterMetrics {
	return &ExporterMetrics{
		scrapingsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "scrapings_total",
			Namespace: namespace,
			Help:      "Scrapings total",
		}, []string{"id"}),
	}
}

func (m *ExporterMetrics) Describe(desc chan<- *prometheus.Desc) {
	m.scrapingsTotal.Describe(desc)
}

func (m *ExporterMetrics) Collect(metric chan<- prometheus.Metric) {
	m.scrapingsTotal.Collect(metric)
}

func (m *ExporterMetrics) IncScrapings(id string) {
	m.scrapingsTotal.WithLabelValues(id).Inc()
}
