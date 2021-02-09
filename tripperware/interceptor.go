package tripperware

import (
	"net/http"

	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/interceptor"
)

// Interceptor tripperware allow multiple req.Body read and allow to set callback before and after roundtrip
func Interceptor(options ...InterceptorOption) httpware.Tripperware {
	config := NewInterceptorConfig(options...)
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
			req.Body = interceptor.NewCopyReadCloser(req.Body)
			config.CallbackBefore(req)
			defer func() {
				config.CallbackAfter(resp, req)
			}()

			return next.RoundTrip(req)
		})
	}
}

type InterceptorConfig struct {
	CallbackBefore func(*http.Request)
	CallbackAfter  func(*http.Response, *http.Request)
}

func (c *InterceptorConfig) apply(options ...InterceptorOption) *InterceptorConfig {
	for _, option := range options {
		option(c)
	}
	return c
}

// NewInterceptorConfig returns a new interceptor configuration with all options applied
func NewInterceptorConfig(options ...InterceptorOption) *InterceptorConfig {
	config := &InterceptorConfig{
		CallbackBefore: func(_ *http.Request) {},
		CallbackAfter:  func(_ *http.Response, _ *http.Request) {},
	}
	return config.apply(options...)
}

// InterceptorOption defines a interceptor tripperware configuration option
type InterceptorOption func(*InterceptorConfig)

// WithAfter will configure CallbackAfter interceptor option
func WithBefore(callbackBefore func(*http.Request)) InterceptorOption {
	return func(config *InterceptorConfig) {
		config.CallbackBefore = callbackBefore
	}
}

// WithAfter will configure CallbackAfter interceptor option
func WithAfter(callbackAfter func(*http.Response, *http.Request)) InterceptorOption {
	return func(config *InterceptorConfig) {
		config.CallbackAfter = callbackAfter
	}
}
