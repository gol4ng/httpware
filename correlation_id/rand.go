package correlation_id

import (
	"math/rand"
	"sync"
)

// LockedSource is a copy of private go lockedSource
// https://github.com/golang/go/blob/master/src/math/rand/rand.go#L374
type LockedSource struct {
	lk  sync.Mutex
	src rand.Source64
}

func (r *LockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *LockedSource) Uint64() (n uint64) {
	r.lk.Lock()
	n = r.src.Uint64()
	r.lk.Unlock()
	return
}

func (r *LockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}

func NewLockedSource(src rand.Source) *LockedSource {
	return &LockedSource{src: src.(rand.Source64)}
}
