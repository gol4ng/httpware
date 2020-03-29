package rate_limit_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gol4ng/httpware/v2/rate_limit"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	assert.Equal(
		t,
		"default error",
		rate_limit.GetConfig().ErrorCallback(fmt.Errorf("default error"), nil).Error(),
	)
}

func TestConfig_Options(t *testing.T) {
	config := rate_limit.GetConfig(
		rate_limit.WithErrorCallback(func(err error, req *http.Request) error {
			return fmt.Errorf("error from callback")
		}),
	)

	assert.Equal(t, "error from callback", config.ErrorCallback(nil, nil).Error())
}
