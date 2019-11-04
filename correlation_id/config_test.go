package correlation_id_test

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware/v2/correlation_id"
)

func TestNewConfig(t *testing.T) {
	defaultRand := rand.New(correlation_id.NewLockedSource(rand.NewSource(1)))
	correlation_id.DefaultIdGenerator = correlation_id.NewRandomIdGenerator(defaultRand)

	for _, expectedId := range []string{"p1LGIehp1s", "uqtCDMLxiD"} {
		assert.Equal(t, expectedId, correlation_id.GetConfig().IdGenerator(nil))
	}
}

func TestConfig_Options(t *testing.T) {
	config := correlation_id.GetConfig(
		correlation_id.WithHeaderName("my-personal-header-name"),
		correlation_id.WithIdGenerator(func(request *http.Request) string {
			return "toto"
		}),
	)
	assert.Equal(t, "my-personal-header-name", config.HeaderName)
	assert.Equal(t, "toto", config.IdGenerator(nil))
}
