package correlation_id_test

import (
	"github.com/gol4ng/httpware/correlation_id"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestNewConfig(t *testing.T) {
	defaultRand := rand.New(correlation_id.NewLockedSource(rand.NewSource(1)))
	correlation_id.DefaultIdGenerator = correlation_id.NewRandomIdGenerator(defaultRand)

	for _, expectedId := range []string{"p1LGIehp1s", "uqtCDMLxiD"} {
		assert.Equal(t, expectedId, correlation_id.NewConfig().IdGenerator(nil))
	}
}
