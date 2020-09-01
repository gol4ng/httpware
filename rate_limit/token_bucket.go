package rate_limit

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	mutex     sync.Mutex
	ticker    *time.Ticker
	done      chan struct{}
	callLimit uint32
	count     uint32
}

func (t *TokenBucket) Allow(_ *http.Request) error {
	t.mutex.Lock()
	res := t.count >= t.callLimit
	t.mutex.Unlock()
	if res {
		return errors.New(RequestLimitReachedErr)
	}

	return nil
}

func (t *TokenBucket) Inc(_ *http.Request) {
	t.mutex.Lock()
	t.count++
	t.mutex.Unlock()
}

func (t *TokenBucket) Dec(_ *http.Request) {}

func (t *TokenBucket) Stop() {
	t.done <- struct{}{}
	t.ticker.Stop()
}

func (t *TokenBucket) start() {
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

func NewTokenBucket(timeBucket time.Duration, callLimit int) *TokenBucket {
	t := &TokenBucket{
		ticker:    time.NewTicker(timeBucket),
		done:      make(chan struct{}),
		callLimit: uint32(callLimit),
	}

	t.start()

	return t
}
