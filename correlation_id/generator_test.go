package correlation_id_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware/v3/correlation_id"
)

func Test_Random(t *testing.T) {
	assert.Equal(t, 10, len(correlation_id.DefaultIdGenerator.Generate(10)))
}

func Test_Random_NewSource(t *testing.T) {
	r := rand.New(correlation_id.NewLockedSource(rand.NewSource(1)))
	rg := correlation_id.NewRandomIdGenerator(r)
	for _, expectedId := range []string{"DHIMG9FpXzp1LGIehp1s", "zAHyfjXUlrGhblT7txWd"} {
		assert.Equal(t, expectedId, rg.Generate(20))
	}
}
