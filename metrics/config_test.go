package metrics_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware/v3/metrics"
)

func TestConfig_Options(t *testing.T) {
	var recorder metrics.Recorder
	config := metrics.NewConfig(
		recorder,
		metrics.WithSplitStatus(true),
		metrics.WithObserveResponseSize(false),
		metrics.WithMeasureInflightRequests(false),
		metrics.WithIdentifierProvider(func(req *http.Request) string {
			return "my-personal-identifier"
		}),
	)
	assert.Equal(t, true, config.SplitStatus)
	assert.Equal(t, false, config.ObserveResponseSize)
	assert.Equal(t, false, config.MeasureInflightRequests)
	assert.Equal(t, "my-personal-identifier", config.IdentifierProvider(nil))
}
