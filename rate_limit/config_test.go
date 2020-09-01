package rate_limit_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gol4ng/httpware/v2/rate_limit"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	_, err := rate_limit.GetConfig().ErrorCallback(nil, fmt.Errorf("default error"))
	assert.NotNil(t, err)
	assert.Equal(
		t,
		"default error",
		err.Error(),
	)
}

func TestConfig_Options(t *testing.T) {
	config := rate_limit.GetConfig(
		rate_limit.WithErrorCallback(func(*http.Request, error) (*http.Response, error) {
			return nil, fmt.Errorf("error from callback")
		}),
	)
	_, err := config.ErrorCallback(nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, "error from callback",err.Error())
}
