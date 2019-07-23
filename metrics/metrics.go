package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Recorder interface {
	RegisterOn(registry prometheus.Registerer) Recorder
	ObserveHTTPRequestDuration(ctx context.Context, id string, duration time.Duration, method, code string)
	ObserveHTTPResponseSize(ctx context.Context, id string, sizeBytes int64, method, code string)
	AddInflightRequests(ctx context.Context, id string, quantity int)
}


