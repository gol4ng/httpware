package rate_limit_test

import (
	"testing"
	"time"

	"github.com/gol4ng/httpware/v2/rate_limit"
	"golang.org/x/time/rate"
)

func BenchmarkLeakyBucket_IsLimitReached(b *testing.B) {
	rl := rate_limit.NewLeakyBucket(1*time.Second, 10)
	go func() {
		for n := 0; n < b.N; n++ {
			rl.Inc()
		}
	}()

	for n := 0; n < b.N; n++ {
		rl.IsLimitReached()
	}
}

func BenchmarkTimeRateLimiter_IsLimitReached(b *testing.B) {
	limiter := rate.NewLimiter(rate.Every(1*time.Second), 10)
	for n := 0; n < b.N; n++ {
		limiter.Allow()
	}
}
