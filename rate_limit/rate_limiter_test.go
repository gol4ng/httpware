package rate_limit_test

import (
	"testing"
	"time"

	"github.com/gol4ng/httpware/v2/rate_limit"
)

func BenchmarkLeakyBucket_IsLimitReached(b *testing.B) {
	rl := rate_limit.NewLeakyBucket(1*time.Second, 10)
	for n := 0; n < b.N; n++ {
		rl.IsLimitReached()
	}
}

func BenchmarkTimeRateLimiter_IsLimitReached(b *testing.B) {
	rl := rate_limit.NewTimeRateLimiter(1*time.Second, 10)
	for n := 0; n < b.N; n++ {
		rl.IsLimitReached()
	}
}
