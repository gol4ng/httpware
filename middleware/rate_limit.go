package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/rate_limit"
)

func RateLimit(rl rate_limit.RateLimiter, options ...RateLimitOption) httpware.Middleware {
	config := NewRateLimitConfig(options...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if err := rl.Allow(req); err != nil {
				if !config.ErrorCallback(err, writer, req) {
					return
				}
			}

			rl.Inc(req)
			defer rl.Dec(req)
			next.ServeHTTP(writer, req)
		})
	}
}

type RateLimitOption func(*RateLimitConfig)

type RateLimitErrorCallback func(err error, writer http.ResponseWriter, req *http.Request) (next bool)

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
		ErrorCallback: DefaultRateLimitErrorCallback,
	}
	return config.apply(options...)
}

func DefaultRateLimitErrorCallback(err error, writer http.ResponseWriter, _ *http.Request) bool {
	http.Error(writer, err.Error(), http.StatusTooManyRequests)
	return false
}

func WithRateLimitErrorCallback(callback RateLimitErrorCallback) RateLimitOption {
	return func(config *RateLimitConfig) {
		config.ErrorCallback = callback
	}
}