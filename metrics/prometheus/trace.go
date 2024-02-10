package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gol4ng/httpware/metrics"
)

func NewInstrumentTrace(config ConfigTrace) *metrics.InstrumentTrace {
	clientDetailsLatencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Subsystem: "http",
			Name:      "details_duration_seconds",
			Help:      "Trace request latency histogram.",
			Buckets:   config.DetailLatencyBuckets,
		},
		[]string{"event"},
	)
	clientDNSLatencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Subsystem: "http",
			Name:      "dns_duration_seconds",
			Help:      "Trace dns latency histogram.",
			Buckets:   config.DNSLatencyBuckets,
		},
		[]string{"event"},
	)
	clientTLSLatencyVec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: config.Namespace,
			Subsystem: "http",
			Name:      "tls_duration_seconds",
			Help:      "Trace tls latency histogram.",
			Buckets:   config.TLSLatencyBuckets,
		},
		[]string{"event"},
	)

	confTrace := &metrics.InstrumentTrace{}

	confTrace.DNSStart = func(ctx context.Context, t time.Duration) {
		clientDNSLatencyVec.WithLabelValues("dns_start")
	}
	confTrace.DNSDone = func(ctx context.Context, t time.Duration) {
		clientDNSLatencyVec.WithLabelValues("dns_done")
	}
	confTrace.ConnectStart = func(ctx context.Context, t time.Duration) {
		clientDetailsLatencyVec.WithLabelValues("connect_start")
	}
	confTrace.ConnectDone = func(ctx context.Context, t time.Duration) {
		clientDetailsLatencyVec.WithLabelValues("connect_done")
	}
	confTrace.TLSHandshakeStart = func(ctx context.Context, t time.Duration) {
		clientTLSLatencyVec.WithLabelValues("tls_handshake_start")
	}
	confTrace.TLSHandshakeDone = func(ctx context.Context, t time.Duration) {
		clientTLSLatencyVec.WithLabelValues("tls_handshake_done")
	}
	confTrace.GotConn = func(ctx context.Context, t time.Duration) {
		clientDetailsLatencyVec.WithLabelValues("got_conn")
	}
	confTrace.PutIdleConn = func(ctx context.Context, t time.Duration) {
		clientDetailsLatencyVec.WithLabelValues("put_idle_conn")
	}
	confTrace.WroteHeaders = func(ctx context.Context, t time.Duration) {
		clientDetailsLatencyVec.WithLabelValues("wrote_headers")
	}
	confTrace.WroteRequest = func(ctx context.Context, t time.Duration) {
		clientDetailsLatencyVec.WithLabelValues("wrote_request")
	}
	confTrace.Got100Continue = func(ctx context.Context, t time.Duration) {
		clientDetailsLatencyVec.WithLabelValues("got_100_continue")
	}
	confTrace.Wait100Continue = func(ctx context.Context, t time.Duration) {
		clientDetailsLatencyVec.WithLabelValues("wait_100_continue")
	}
	confTrace.GotFirstResponseByte = func(ctx context.Context, t time.Duration) {
		clientDetailsLatencyVec.WithLabelValues("got_first_response_byte")
	}

	prometheus.MustRegister(clientDetailsLatencyVec, clientDNSLatencyVec, clientTLSLatencyVec)

	return confTrace
}
