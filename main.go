package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/nightfury1204/prometheus-remote-metric-writer/metrics"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	metricsConf := metrics.NewMetricsExporterConfigs()

	var rootCmd = &cobra.Command{
		Use:               "metrics-writer [command]",
		Short:             `Prometheus metrics writer`,
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := metricsConf.Validate(); err != nil {
				return err
			}

			metricsExporter, err := metrics.NewMetricsExporter(metricsConf, prometheus.NewRegistry())
			if err != nil {
				return errors.Wrap(err, "failed to create client for metrics exporter")
			}

			stopCh := make(chan struct{})
			if err := metricsExporter.Run(stopCh); err != nil {
				return err
			}
			<-stopCh
			return nil
		},
	}
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	flag.CommandLine.Parse([]string{})
	metricsConf.AddFlags(rootCmd.PersistentFlags())
	return rootCmd
}

func main() {
	rootCmd := NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		glog.Fatal(err)
	}
}
