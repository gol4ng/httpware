package rate_limit

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	mux sync.Mutex

	timeBucket time.Duration
	ticker *time.Ticker
	done chan bool
	isStart bool
	callLimit int
	count int
}

func (t *LeakyBucket) IsLimitReached() bool {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.count >= t.callLimit
}

func (t *LeakyBucket) Inc() {
	t.mux.Lock()
	defer t.mux.Unlock()

	t.count++
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
				t.mux.Lock()
				t.count = 0
				t.mux.Unlock()
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
		done: make(chan bool),
		callLimit: callLimit,
	}
}
