package rate_limit

import (
	"net/http"
)

type RateLimiter interface {
	IsLimitReached(req *http.Request) bool
	Inc(req *http.Request)
	Dec(req *http.Request)
}
