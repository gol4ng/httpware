package rate_limit_test

import (
	"testing"
	"time"

	"github.com/gol4ng/httpware/v2/rate_limit"
	"github.com/stretchr/testify/assert"
)

func TestTimeRateLimiter_IsLimitReached(t *testing.T) {
	rl := rate_limit.NewTimeRateLimiter(1*time.Millisecond, 1)
	assert.Equal(t, false, rl.IsLimitReached())
	assert.Equal(t, true, rl.IsLimitReached())

	time.Sleep(2*time.Millisecond)
	assert.Equal(t, false, rl.IsLimitReached())
}
