package correlation_id

import (
	"math/rand"
	"time"
	"unsafe"
)

//https://stackoverflow.com/a/31832326
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

type RandomIdGenerator struct {
	r *rand.Rand
}

func (rg *RandomIdGenerator) Generate(length int) string {
	b := make([]byte, length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, rg.r.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rg.r.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

var DefaultIdGenerator = NewRandomIdGenerator(
	rand.New(
		NewLockedSource(
			rand.NewSource(
				time.Now().UTC().UnixNano(),
			),
		),
	),
)

func NewRandomIdGenerator(rand *rand.Rand) *RandomIdGenerator {
	return &RandomIdGenerator{
		r: rand,
	}
}
