package rate_limit_test

import (
	"testing"
	"time"

	"github.com/gol4ng/httpware/v2/rate_limit"
	"github.com/stretchr/testify/assert"
)

func TestLeakyBucket_IsLimitReached(t *testing.T) {
	rl := rate_limit.NewLeakyBucket(1 * time.Millisecond, 1)
	defer rl.Stop()

	assert.Equal(t, false, rl.IsLimitReached())
	rl.Inc()

	assert.Equal(t, true, rl.IsLimitReached())
	rl.Inc()

	time.Sleep(2 * time.Millisecond)
	assert.Equal(t, false, rl.IsLimitReached())
}

func TestLeakyBucket_Race(t *testing.T) {
	rl := rate_limit.NewLeakyBucket(1 * time.Millisecond, 1)
	for i := 0; i <= 1000; i++ {
		go func() {
			rl.IsLimitReached()
		}()
	}
}
