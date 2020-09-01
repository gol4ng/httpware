package rate_limit_test

import (
	"testing"
	"time"

	"github.com/gol4ng/httpware/v2/rate_limit"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_Allow(t *testing.T) {
	rl := rate_limit.NewTokenBucket(1 * time.Millisecond, 1)
	defer rl.Stop()

	assert.NoError(t, rl.Allow(nil))
	rl.Inc(nil)

	assert.EqualError(t, rl.Allow(nil), "request limit reached")
	rl.Inc(nil)

	time.Sleep(2 * time.Millisecond)
	assert.NoError(t, rl.Allow(nil))
}
