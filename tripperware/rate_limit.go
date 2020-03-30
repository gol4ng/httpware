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
		return httpware.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if rateLimiter.IsLimitReached() {
				return nil, config.ErrorCallback(errors.New(rate_limit.RequestLimitReachedErr), req)
			}

			rateLimiter.Inc()
			return next.RoundTrip(req)
		})
	}
}
