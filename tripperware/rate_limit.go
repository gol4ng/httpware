package tripperware

import (
	"errors"
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/rate_limit"
)

func RateLimit(rateLimiter rate_limit.RateLimiter, options ...rate_limit.Option) httpware.Tripperware {
	config := rate_limit.GetConfig(options...)

	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			if rateLimiter.IsLimitReached(request) {
				return config.ErrorCallback(request, errors.New(rate_limit.RequestLimitReachedErr))
			}

			rateLimiter.Inc(request)
			defer rateLimiter.Dec(request)
			return next.RoundTrip(request)
		})
	}
}
