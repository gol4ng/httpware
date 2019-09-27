package request_id_test

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"

	"github.com/gol4ng/httpware/request_id"
)

func Test_LockedSource_Int63(t *testing.T) {
	s := request_id.NewLockedSource(rand.NewSource(1))

	for _, v := range []int64{5577006791947779410, 8674665223082153551} {
		assert.Equal(t, v, s.Int63())
	}
}

func Test_LockedSource_Uint64(t *testing.T) {
	s := request_id.NewLockedSource(rand.NewSource(1))

	for _, v := range []uint64{0x4d65822107fcfd52, 0x78629a0f5f3f164f} {
		assert.Equal(t, v, s.Uint64())
	}
}

func Test_LockedSource_Seed(t *testing.T) {
	s := request_id.NewLockedSource(rand.NewSource(1))

	for _, v := range []uint64{0x4d65822107fcfd52, 0x78629a0f5f3f164f} {
		assert.Equal(t, v, s.Uint64())
	}
	s.Seed(123)
	for _, v := range []uint64{0x4a68998bed5c40f1, 0x835b51599210f9ba} {
		assert.Equal(t, v, s.Uint64())
	}
}
