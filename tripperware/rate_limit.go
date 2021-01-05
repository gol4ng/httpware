package tripperware

import (
	"net/http"

	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/rate_limit"
)

func RateLimit(rateLimiter rate_limit.RateLimiter, options ...RateLimitOption) httpware.Tripperware {
	config := NewRateLimitConfig(options...)

	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			if limitErr := rateLimiter.Allow(request); limitErr != nil {
				if res, err := config.ErrorCallback(request, limitErr); err != nil {
					return res, err
				}
			}

			rateLimiter.Inc(request)
			defer rateLimiter.Dec(request)
			return next.RoundTrip(request)
		})
	}
}


type RateLimitErrorCallback func(request *http.Request, limitErr error) (response *http.Response, err error)

type RateLimitOption func(*RateLimitConfig)

type RateLimitConfig struct {
	ErrorCallback RateLimitErrorCallback
}

func (c *RateLimitConfig) apply(options ...RateLimitOption) *RateLimitConfig {
	for _, option := range options {
		option(c)
	}
	return c
}

func NewRateLimitConfig(options ...RateLimitOption) *RateLimitConfig {
	config := &RateLimitConfig{
		ErrorCallback: DefaultRateLimitErrorCallback(),
	}
	return config.apply(options...)
}

func DefaultRateLimitErrorCallback() RateLimitErrorCallback {
	return func(_ *http.Request, err error) (*http.Response, error) {
		return nil, err
	}
}

func WithRateLimitErrorCallback(callback RateLimitErrorCallback) RateLimitOption {
	return func(config *RateLimitConfig) {
		config.ErrorCallback = callback
	}
}
