package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type TestMetricsCollector struct {
	metrics *prometheus.Desc
}

func NewTestMetricsCollector(id string) *TestMetricsCollector {
	return &TestMetricsCollector{
		metrics: prometheus.NewDesc(
			prometheus.BuildFQName("test_metrics", "", id),
			"test metrics",
			nil, prometheus.Labels{
				"who" : "metric-exporter",
			},
		),
	}
}

func (c *TestMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.metrics
}

func (c *TestMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.metrics, prometheus.GaugeValue, 1)
}
