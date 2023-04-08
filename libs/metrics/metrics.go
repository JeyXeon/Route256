package metrics

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ServerRequestsCounter       prometheus.Counter
	ServerResponseCounter       *prometheus.CounterVec
	ServerHistogramResponseTime *prometheus.HistogramVec
	ClientHistogramResponseTime *prometheus.HistogramVec
)

func Init() {
	reg := prometheus.NewRegistry()
	ServerRequestsCounter = promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "server_requests_total",
	})
	ServerResponseCounter = promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "server_responses_total",
	},
		[]string{"status"},
	)
	ServerHistogramResponseTime = promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "server_histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status"},
	)
	ClientHistogramResponseTime = promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "client_histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status"},
	)
}

func ServeMetrics(port string, logger *zap.Logger) {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatal("failed to listen and serve metrics", zap.Error(err))
	}
}
