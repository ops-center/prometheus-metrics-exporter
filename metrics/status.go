package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName("", "", "up"),
		"Is metrics exporter is running.",
		nil, prometheus.Labels{
			"version": "dev",
			"who" : "metric-exporter",
		},
	)
)

type HealthCollector struct {
}

func NewHealthCollector() *HealthCollector {
	return &HealthCollector{}
}

func (c *HealthCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

func (c *HealthCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)
}
