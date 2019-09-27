package request_id_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware/request_id"
)

func Test_Random(t *testing.T) {
	assert.Equal(t, 10, len(request_id.DefaultIdGenerator.Generate(nil)))
}

func Test_Random_NewSource(t *testing.T) {
	r := rand.New(request_id.NewLockedSource(rand.NewSource(1)))
	rg := request_id.NewRandomIdGenerator(r, 20)
	for _, expectedId := range []string{"DHIMG9FpXzp1LGIehp1s", "zAHyfjXUlrGhblT7txWd"} {
		assert.Equal(t, expectedId, rg.Generate(nil))
	}
}
