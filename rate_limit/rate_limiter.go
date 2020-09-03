package rate_limit

import (
	"net/http"
)

type RateLimiter interface {
	Allow(req *http.Request) error
	Inc(req *http.Request)
	Dec(req *http.Request)
}
