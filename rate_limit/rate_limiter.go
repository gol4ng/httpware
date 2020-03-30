package rate_limit

type RateLimiter interface {
	IsLimitReached() bool
	Inc()
}
