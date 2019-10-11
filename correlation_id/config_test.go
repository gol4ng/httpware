package correlation_id_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware/correlation_id"
)

func TestNewConfig(t *testing.T) {
	defaultRand := rand.New(correlation_id.NewLockedSource(rand.NewSource(1)))
	correlation_id.DefaultIdGenerator = correlation_id.NewRandomIdGenerator(defaultRand)

	for _, expectedId := range []string{"p1LGIehp1s", "uqtCDMLxiD"} {
		assert.Equal(t, expectedId, correlation_id.NewConfig().IdGenerator(nil))
	}
}
