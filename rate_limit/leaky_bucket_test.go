package rate_limit_test

import (
	"testing"
	"time"

	"github.com/gol4ng/httpware/v2/rate_limit"
	"github.com/stretchr/testify/assert"
)

func TestLeakyBucket_IsLimitReached(t *testing.T) {
	rl := rate_limit.NewLeakyBucket(1*time.Millisecond, 1)
	rl.Start()
	defer rl.Stop()

	rl.Inc()
	assert.Equal(t, true, rl.IsLimitReached())

	rl.Inc()
	assert.Equal(t, true, rl.IsLimitReached())

	time.Sleep(2*time.Millisecond)
	assert.Equal(t, false, rl.IsLimitReached())
}

func TestLeakyBucket_StartStop(t *testing.T) {
	rl := rate_limit.NewLeakyBucket(1*time.Millisecond, 1)
	rl.Start()
	rl.Stop()

	rl.Inc()
	rl.Inc()
	time.Sleep(2*time.Millisecond)
	assert.Equal(t, true, rl.IsLimitReached())

	rl.Start()
	rl.Inc()
	rl.Inc()
	time.Sleep(2*time.Millisecond)
	assert.Equal(t, false, rl.IsLimitReached())
	rl.Stop()
}

func TestLeakyBucket_Race(t *testing.T) {
	rl := rate_limit.NewLeakyBucket(1*time.Millisecond, 1)
	for i := 0; i <= 1000; i++ {
		go func() {
			rl.Inc()
		}()
	}

	for i := 0; i <= 1000; i++ {
		go func() {
			rl.IsLimitReached()
		}()
	}
}
