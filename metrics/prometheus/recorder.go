package prometheus

import (
	"context"
	"time"

	"github.com/gol4ng/httpware/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

type Recorder struct {
	httpRequestDurHistogram   *prometheus.HistogramVec
	httpResponseSizeHistogram *prometheus.HistogramVec
	httpRequestsInflight      *prometheus.GaugeVec
}

func (r *Recorder) RegisterOn(registry prometheus.Registerer) metrics.Recorder {
	if registry == nil {
		registry = prometheus.DefaultRegisterer
	}

	registry.MustRegister(
		r.httpRequestDurHistogram,
		r.httpResponseSizeHistogram,
		r.httpRequestsInflight,
	)
	return r
}

func (r *Recorder) ObserveHTTPRequestDuration(_ context.Context, id string, duration time.Duration, method, code string) {
	r.httpRequestDurHistogram.WithLabelValues(id, method, code).Observe(duration.Seconds())
}

func (r *Recorder) ObserveHTTPResponseSize(_ context.Context, id string, responseSize int64, method, code string) {
	r.httpResponseSizeHistogram.WithLabelValues(id, method, code).Observe(float64(responseSize))
}

func (r *Recorder) AddInflightRequests(_ context.Context, id string, quantity int) {
	r.httpRequestsInflight.WithLabelValues(id).Add(float64(quantity))
}

func NewRecorder(config Config) *Recorder {
	config.defaults()

	return &Recorder{
		httpRequestDurHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: config.Prefix,
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "The latency of the HTTP requests (in seconds).",
			Buckets:   config.DurationBuckets,
		}, []string{config.IdentifierLabel, config.MethodLabel, config.StatusCodeLabel}),
		httpResponseSizeHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: config.Prefix,
			Subsystem: "http",
			Name:      "response_size_bytes",
			Help:      "The size of the HTTP responses (in bytes).",
			Buckets:   config.SizeBuckets,
		}, []string{config.IdentifierLabel, config.MethodLabel, config.StatusCodeLabel}),
		httpRequestsInflight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: config.Prefix,
			Subsystem: "http",
			Name:      "requests_inflight",
			Help:      "The number of inflight requests being handled at the same time.",
		}, []string{config.IdentifierLabel}),
	}
}
