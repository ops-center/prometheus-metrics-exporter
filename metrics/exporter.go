package metrics

import (
	"github.com/golang/glog"
	"github.com/prometheus/prometheus/prompb"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	prom_config "github.com/prometheus/common/config"
	"github.com/spf13/pflag"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	defaultInterval = time.Second * 30
	defaultTimeout  = time.Minute * 3
)

type MetricsExporterConfigs struct {
	// The id of the metrics exporter
	Id string

	// The address where metrics will be sent
	Addr string

	// Interval at which metrics data will be sent
	Interval time.Duration

	// Metrics write timeout
	WriteTimeout time.Duration

	// The CA cert to use for the targets.
	CAFile string

	// The client cert file for the targets.
	CertFile string

	// The client key file for the targets.
	KeyFile string

	// Used to verify the hostname for the targets.
	ServerName string

	// Disable target certificate validation.
	InsecureSkipVerify bool

	// License to use for authentication
	License string
}

func NewMetricsExporterConfigs() *MetricsExporterConfigs {
	return &MetricsExporterConfigs{}
}

func (m *MetricsExporterConfigs) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&m.Id, "metrics-exporter.id", m.Id, "The id for metrics exporter")
	fs.StringVar(&m.Addr, "metrics-exporter.url", m.Addr, "The address of metrics storage where metrics data will be sent")
	fs.DurationVar(&m.WriteTimeout, "metrics-exporter.write-timeout", defaultTimeout, "Specifies the metrics write timeout")
	fs.DurationVar(&m.Interval, "metrics-exporter.interval", defaultInterval, "Specifies the interval at which metrics data will be sent")
	fs.StringVar(&m.CAFile, "metrics-exporter.ca-cert-file", m.CAFile, "The path of the CA cert to use for the remote metric storage.")
	fs.StringVar(&m.CertFile, "metrics-exporter.client-cert-file", m.CertFile, "The path of the client cert to use for communicating with the remote metric storage.")
	fs.StringVar(&m.KeyFile, "metrics-exporter.client-key-file", m.KeyFile, "The path of the client key to use for communicating with the remote metric storage.")
	fs.StringVar(&m.ServerName, "metrics-exporter.server-name", m.ServerName, "The server name which will be used to verify metrics storage.")
	fs.BoolVar(&m.InsecureSkipVerify, "metrics-exporter.insecure-skip-verify", m.InsecureSkipVerify, "To skip tls verification when communicating with the remote metric storage.")
	fs.StringVar(&m.License, "metrics-exporter.license", m.License, "License to use for authentication")
}

func (m *MetricsExporterConfigs) Validate() error {
	if m.Id == "" {
		return errors.New("metrics-exporter.id must non-empty")
	}
	if m.Addr == "" {
		return errors.New("metrics-exporter.url must non-empty")
	}
	if m.WriteTimeout < time.Second*5 {
		return errors.New("metrics-exporter.write-timeout must be greater than 5s")
	}
	if m.Interval < time.Second*5 {
		return errors.New("metrics-exporter.write-timeout must be greater than 5s")
	}
	return nil
}

type MetricsExporter struct {
	Config *MetricsExporterConfigs

	// Prometheus registry
	PromRegistry *prometheus.Registry
}

// registry is nil, it will create new one
func NewMetricsExporter(c *MetricsExporterConfigs, registry *prometheus.Registry) (*MetricsExporter, error) {
	if registry == nil {
		registry = prometheus.NewRegistry()
	}

	// TODO: another function?
	// TODO: remove later
	registry.MustRegister(NewHealthCollector(),
		NewTestMetricsCollector(c.Id))

	return &MetricsExporter{
		Config:       c,
		PromRegistry: registry,
	}, nil
}

func (m *MetricsExporter) Register(cs ...prometheus.Collector) {
	m.PromRegistry.MustRegister(cs...)
}

// non-blocking
func (m *MetricsExporter) Run(stopCh <-chan struct{}) error {
	httpConf := prom_config.HTTPClientConfig{
		TLSConfig: prom_config.TLSConfig{
			CAFile:             m.Config.CAFile,
			CertFile:           m.Config.CertFile,
			KeyFile:            m.Config.KeyFile,
			ServerName:         m.Config.ServerName,
			InsecureSkipVerify: m.Config.InsecureSkipVerify,
		},
	}

	if m.Config.License != "" {
		// set the license as bearer token
		httpConf.BearerToken = prom_config.Secret(m.Config.License)

		jwtToken, err := jwt.ParseSigned(m.Config.License)
		if err != nil {
			return err
		}

		data := &License{}
		err = jwtToken.UnsafeClaimsWithoutVerification(data)
		if err != nil {
			return err
		}

		glog.Info("Overwriting client id from license", "client_id", data.Subject)
		m.Config.Id = data.Subject

	} else {
		glog.Warning("license is not provided")
	}

	cl, err := NewRemoteClient(m.Config.Addr, httpConf, m.Config.WriteTimeout)
	if err != nil {
		return errors.Wrap(err, "failed to create metrics storage remote client")
	}

	// TODO: all extra labels in here
	var extraLabels []prompb.Label
	extraLabels = append(extraLabels, GetLabels()...)
	extraLabels = append(extraLabels, prompb.Label{
		Name: "client_id",
		Value: m.Config.Id,
	})

	rw, err := NewRemoteWriter(cl, m.PromRegistry, m.Config.Interval, extraLabels)
	if err != nil {
		return errors.Wrap(err, "failed to create remote writer for metrics")
	}

	go rw.Run(stopCh)

	return nil
}
