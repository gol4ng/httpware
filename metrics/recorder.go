package metrics

import (
	"context"
	"time"
)

type Recorder interface {
	ObserveHTTPRequestDuration(ctx context.Context, id string, duration time.Duration, method, code string)
	ObserveHTTPResponseSize(ctx context.Context, id string, responseSize int64, method, code string)
	AddInflightRequests(ctx context.Context, id string, quantity int)
}
