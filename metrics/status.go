package metrics

import (
	v "github.com/appscode/go/version"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName("operator", "kubevault", "up"),
		"Is kubevault operator is running.",
		nil, prometheus.Labels{
			"maintainer": "appscode",
			"version":    v.Version.Version,
		},
	)
)

type OperatorHealthCollector struct {
}

func NewOperatorHealthCollector() *OperatorHealthCollector {
	return &OperatorHealthCollector{}
}

func (c *OperatorHealthCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

func (c *OperatorHealthCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)
}
