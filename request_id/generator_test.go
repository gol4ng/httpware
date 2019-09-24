package request_id_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware/request_id"
)

func Test_Random(t *testing.T) {
	assert.Equal(t, 10, len(request_id.RandomIdGenerator(nil)))
}

func Test_Random_NewSource(t *testing.T) {
	request_id.Rand = rand.New(rand.NewSource(1))
	for _, expectedId := range []string{"p1LGIehp1s", "uqtCDMLxiD"} {
		assert.Equal(t, expectedId, request_id.RandomIdGenerator(nil))
	}
}
