package rate_limit_test

import (
	"testing"
	"time"

	"github.com/gol4ng/httpware/v2/rate_limit"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_IsLimitReached(t *testing.T) {
	rl := rate_limit.NewTokenBucket(1 * time.Millisecond, 1)
	defer rl.Stop()

	assert.Equal(t, false, rl.IsLimitReached(nil))
	rl.Inc(nil)

	assert.Equal(t, true, rl.IsLimitReached(nil))
	rl.Inc(nil)

	time.Sleep(2 * time.Millisecond)
	assert.Equal(t, false, rl.IsLimitReached(nil))
}
