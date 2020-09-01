package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/rate_limit"
)

func RateLimit(rl rate_limit.RateLimiter) httpware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if rl.IsLimitReached(req) {
				http.Error(writer, rate_limit.RequestLimitReachedErr, http.StatusTooManyRequests)
				return
			}

			rl.Inc(req)
			defer rl.Dec(req)
			next.ServeHTTP(writer, req)
		})
	}
}
