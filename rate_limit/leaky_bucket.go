package rate_limit

import (
	"sync/atomic"
	"time"
)

type LeakyBucket struct {
	timeBucket time.Duration
	ticker     *time.Ticker
	done       chan bool
	isStart    bool
	callLimit  uint64
	count      uint64
}

func (t *LeakyBucket) IsLimitReached() bool {
	if atomic.LoadUint64(&t.count) >= t.callLimit {
		return true
	}

	atomic.AddUint64(&t.count, 1)
	return false
}

func (t *LeakyBucket) Start() {
	if t.isStart {
		return
	}

	t.ticker = time.NewTicker(t.timeBucket)
	t.isStart = true
	go func() {
		for {
			select {
			case <-t.done:
				return
			case <-t.ticker.C:
				atomic.StoreUint64(&t.count, 0)
			}
		}
	}()
}

func (t *LeakyBucket) Stop() {
	t.done <- true
	t.ticker.Stop()
	t.isStart = false
}

func NewLeakyBucket(timeBucket time.Duration, callLimit int) *LeakyBucket {
	return &LeakyBucket{
		timeBucket: timeBucket,
		done:       make(chan bool),
		callLimit:  uint64(callLimit),
	}
}
