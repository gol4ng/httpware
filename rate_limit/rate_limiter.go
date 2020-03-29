package rate_limit

type RateLimiter interface {
	Start()
	Stop()
	Inc()
	IsLimitReached() bool
}
