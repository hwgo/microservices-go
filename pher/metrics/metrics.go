package metrics

import (
	jm "github.com/uber/jaeger-lib/metrics"
	"github.com/uber/jaeger-lib/metrics/go-kit"
	"github.com/uber/jaeger-lib/metrics/go-kit/expvar"
	jprom "github.com/uber/jaeger-lib/metrics/prometheus"
)

var (
	defaultBackend string
	metricsFactory jm.Factory
)

func init() {
	// defaultBackend = "expvar"
	defaultBackend = "prometheus"
	initMetrics(defaultBackend)
}

func initMetrics(metricsBackend string) {
	if metricsBackend == "expvar" {
		metricsFactory = xkit.Wrap("", expvar.NewFactory(10)) // 10 buckets for histograms
		// logger.Info("Using expvar as metrics backend")
	} else if metricsBackend == "prometheus" {
		metricsFactory = jprom.New()
		// logger.Info("Using Prometheus as metrics backend")
	} else {
		// logger.Fatal("unsupported metrics backend " + metricsBackend)
	}
}

func DefaultMetricsFactory() jm.Factory {
	return metricsFactory
}

func Namespace(name string, tags map[string]string) jm.Factory {
	return metricsFactory.Namespace(name, tags)
}
