package rate_limit_test

import (
	"testing"
	"time"

	"github.com/gol4ng/httpware/v2/rate_limit"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestLeakyBucket_IsLimitReached(t *testing.T) {
	rl := rate_limit.NewLeakyBucket(1*time.Millisecond, 1)
	rl.Start()
	defer rl.Stop()

	assert.Equal(t, false, rl.IsLimitReached())
	assert.Equal(t, true, rl.IsLimitReached())

	time.Sleep(2 * time.Millisecond)
	assert.Equal(t, false, rl.IsLimitReached())
}

func TestLeakyBucket_StartStop(t *testing.T) {
	rl := rate_limit.NewLeakyBucket(1*time.Millisecond, 1)
	rl.Start()
	rl.Stop()

	rl.IsLimitReached()
	rl.IsLimitReached()
	time.Sleep(2 * time.Millisecond)
	assert.Equal(t, true, rl.IsLimitReached())

	rl.Start()
	rl.IsLimitReached()
	rl.IsLimitReached()
	time.Sleep(2 * time.Millisecond)
	assert.Equal(t, false, rl.IsLimitReached())
	rl.Stop()
}

func TestLeakyBucket_Race(t *testing.T) {
	rl := rate_limit.NewLeakyBucket(1*time.Millisecond, 1)
	for i := 0; i <= 1000; i++ {
		go func() {
			rl.IsLimitReached()
		}()
	}
}

func TestTimeRate(t *testing.T) {
	l := rate.NewLimiter(rate.Every(1*time.Second), 1)
	assert.Equal(t, true, l.Allow())
	assert.Equal(t, false, l.Allow())

	time.Sleep(2 * time.Second)
	assert.Equal(t, true, l.Allow())
}
