package rate_limit

import (
	"time"

	"golang.org/x/time/rate"
)

type TimeRateLimiter struct {
	limiter *rate.Limiter
}

func (rl *TimeRateLimiter) IsLimitReached() bool {
	return !rl.limiter.Allow()
}

func NewTimeRateLimiter(timeBucket time.Duration, callLimit int) *TimeRateLimiter {
	return &TimeRateLimiter{rate.NewLimiter(rate.Every(timeBucket), callLimit)}
}
