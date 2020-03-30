package rate_limit

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	mutex     sync.Mutex
	ticker    *time.Ticker
	done      chan struct{}
	callLimit uint32
	count     uint32
}

func (t *LeakyBucket) IsLimitReached() bool {
	t.mutex.Lock()
	res := t.count >= t.callLimit
	t.mutex.Unlock()
	return res
}

func (t *LeakyBucket) Inc() {
	t.mutex.Lock()
	t.count++
	t.mutex.Unlock()
}

func (t *LeakyBucket) Stop() {
	t.done <- struct{}{}
	t.ticker.Stop()
}

func (t *LeakyBucket) start() {
	go func() {
		for {
			select {
			case <-t.done:
				return
			case <-t.ticker.C:
				t.mutex.Lock()
				t.count = 0
				t.mutex.Unlock()
			}
		}
	}()
}

func NewLeakyBucket(timeBucket time.Duration, callLimit int) *LeakyBucket {
	t := &LeakyBucket{
		ticker:    time.NewTicker(timeBucket),
		done:      make(chan struct{}),
		callLimit: uint32(callLimit),
	}

	t.start()

	return t
}
