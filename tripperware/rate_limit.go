package tripperware

import (
	"fmt"
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/rate_limit"
)

const RequestLimitReachedErr = "request limit reached"

type RateLimitOptions struct {
	errorCallback func(err error, req *http.Request) error
}

type RateLimitOption func(*RateLimitOptions)

func NopErrorCallback() func(err error, req *http.Request) error {
	return func(err error, req *http.Request) error {
		return err
	}
}

func WithErrorCallback(callback func(err error, req *http.Request) error) RateLimitOption {
	return func(args *RateLimitOptions) {
		args.errorCallback = callback
	}
}

func RateLimit(rateLimiter rate_limit.RateLimiter, options ...RateLimitOption) httpware.Tripperware {
	args := &RateLimitOptions{
		errorCallback: NopErrorCallback(),
	}

	for _, setter := range options {
		setter(args)
	}

	rateLimiter.Start()
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if rateLimiter.IsLimitReached() {
				return nil, args.errorCallback(fmt.Errorf(RequestLimitReachedErr), req)
			}

			rateLimiter.Inc()
			return next.RoundTrip(req)
		})
	}
}
