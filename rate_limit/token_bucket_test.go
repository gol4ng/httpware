package rate_limit_test

import (
	"testing"
	"time"

	"github.com/gol4ng/httpware/v3/rate_limit"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_Allow(t *testing.T) {
	limiter := rate_limit.NewTokenBucket(1 * time.Millisecond, 1)
	defer limiter.Stop()

	assert.NoError(t, limiter.Allow(nil))
	limiter.Inc(nil)

	assert.EqualError(t, limiter.Allow(nil), "request limit reached")
	limiter.Inc(nil)

	time.Sleep(2 * time.Millisecond)
	assert.NoError(t, limiter.Allow(nil))
}
