package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prebid/prebid-server/config"
	metricsconfig "github.com/prebid/prebid-server/pbsmetrics/config"
	prometheusMetrics "github.com/prebid/prebid-server/pbsmetrics/prometheus"
)

func newPrometheusServer(cfg *config.Configuration, metrics *metricsconfig.DetailedMetricsEngine) *http.Server {
	var proMetrics *prometheusMetrics.Metrics

	proMetrics = metrics.PrometheusMetrics

	if proMetrics == nil {
		glog.Fatal("Prometheus metrics configured, but a Prometheus metrics engine was not found. Cannot set up a Prometheus listener.")
	}
	return &http.Server{
		Addr: cfg.Host + ":" + strconv.Itoa(cfg.Metrics.Prometheus.Port),
		Handler: promhttp.HandlerFor(proMetrics.Registry, promhttp.HandlerOpts{
			ErrorLog:            loggerForPrometheus{},
			MaxRequestsInFlight: 5,
			Timeout:             5 * time.Second,
		}),
	}
}

type loggerForPrometheus struct{}

func (loggerForPrometheus) Println(v ...interface{}) {
	glog.Warningln(v...)
}
