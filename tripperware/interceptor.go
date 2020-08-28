package tripperware

import (
	"net/http"

	"github.com/gol4ng/httpware/v3"
	"github.com/gol4ng/httpware/v3/interceptor"
)

// Interceptor tripperware allow multiple req.Body read and allow to set callback before and after roundtrip
func Interceptor(options ...Option) httpware.Tripperware {
	config := NewConfig(options...)
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

type Config struct {
	CallbackBefore func(*http.Request)
	CallbackAfter  func(*http.Response, *http.Request)
}

func (c *Config) apply(options ...Option) *Config {
	for _, option := range options {
		option(c)
	}
	return c
}

// NewConfig returns a new interceptor configuration with all options applied
func NewConfig(options ...Option) *Config {
	config := &Config{
		CallbackBefore: func(_ *http.Request) {},
		CallbackAfter:  func(_ *http.Response, _ *http.Request) {},
	}
	return config.apply(options...)
}

// Option defines a interceptor tripperware configuration option
type Option func(*Config)

// WithAfter will configure CallbackAfter interceptor option
func WithBefore(callbackBefore func(*http.Request)) Option {
	return func(config *Config) {
		config.CallbackBefore = callbackBefore
	}
}

// WithAfter will configure CallbackAfter interceptor option
func WithAfter(callbackAfter func(*http.Response, *http.Request)) Option {
	return func(config *Config) {
		config.CallbackAfter = callbackAfter
	}
}
