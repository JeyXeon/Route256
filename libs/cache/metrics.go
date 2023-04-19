package cache

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CacheHitCount              prometheus.Counter
	CacheRequestsTotal         prometheus.Counter
	CacheHistogramResponseTime *prometheus.HistogramVec
)

func InitCacheMetrics() {
	CacheHitCount = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "cache_hits_total",
	})

	CacheRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "cache_requests_total",
	})

	CacheHistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "cache_histogram_response_time_seconds",
		Buckets: prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status"},
	)
}
