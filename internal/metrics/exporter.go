package metrics

import "github.com/prometheus/client_golang/prometheus"

type ExporterMetrics struct {
	info           *prometheus.GaugeVec
	scrapingsTotal *prometheus.CounterVec
}

func NewExporterMetrics(namespace string) *ExporterMetrics {
	return &ExporterMetrics{
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:      "info",
			Namespace: namespace,
			Help:      "Exporter info",
		}, []string{"version", "build_date"}),
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

func (m *ExporterMetrics) SetInfo(version, buildDate string) {
	m.info.WithLabelValues(version, buildDate).Set(1)
}
